package main

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"google.golang.org/api/iterator"
)

// Recieves Botrequest, saves message to firebase, sends to recipient, and returns reponse
func (bot *Bot) HandleBusinessInput(reqObj BusinessRequest) BusinessResponse {
	messageRef := bot.Client.Collection(Businesses).Doc(reqObj.Business.Id).Collection(Messages).NewDoc()

	timeInMil := time.Now().UnixNano() / 1000000

	message := Message{
		Content:          reqObj.Message,
		IsBusinessSender: true,
		TimeSent:         timeInMil,
		HasBusinessRead:  true,
		DidBotCreate:     false,
		RecipientId:      reqObj.Recipient.Id}

	_, err := messageRef.Set(bot.Ctx, message)

	personRef := bot.Client.Collection(Businesses).Doc(reqObj.BusinessId).Collection(Recipients).Doc(reqObj.Recipient.Id)
	personRef.Update(bot.Ctx, []firestore.Update{
		{Path: RecentMessage, Value: message},
	})

	if err != nil {
		log.Println(err)
	}

	message.Id = messageRef.ID

	smsRequest := SMSRequest{
		To:   reqObj.Recipient.Contact,
		From: reqObj.Business.PhoneNumber,
		Body: reqObj.Message}

	bot.SmsClient.SendSMS(smsRequest)

	return BusinessResponse{}
}


func (bot *Bot) HandleOutsideInput(reqObj OutsideRequest) OutsideResponse {

  // Need to save the new message to firebase
  err := bot.saveMessage(reqObj.Business, reqObj.Recipient, reqObj.Message)

  if err != nil {
    log.Println(err)
  }

  businessId := reqObj.Business.Id

  // Need to check if a recipient was found, and if not create one, and if so update the recent message
	if reqObj.Recipient.Id == "" {
		reqObj.Recipient.RecentMessage = reqObj.Message

		personRef, _, err := bot.Client.Collection(Businesses).Doc(businessId).Collection(Recipients).Add(bot.Ctx, reqObj.Recipient)
		reqObj.Recipient.Id = personRef.ID

    if err != nil {
      log.Println(err)
    }
	} else {
		personRef := bot.Client.Collection(Businesses).Doc(businessId).Collection(Recipients).Doc(reqObj.Recipient.Id)
		personRef.Update(ctx, []firestore.Update{
			{Path: RecentMessage, Value: reqObj.Message},
		})
	}

  // Send a http request that will be handled in the textual_input_channel
  // The body is the OutsideRequest object
	body, err := json.Marshal(reqObj)

  if err != nil {
    log.Println(err)
  }


	rasaUrl := fmt.Sprintf("http://localhost:5005/webhooks/textual/webhook")
	req, err := http.NewRequest("POST", rasaUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Println(err)
	}

	http.DefaultClient.Do(req)

  return OutsideResponse{}
}

func (bot *Bot) HandleAction(req *RasaRequest) (*RasaResponse, error) {
  resp := NewRasaResponse()

	action := req.NextAction
  log.Println(action)
		switch action {
		case ACTION_START_ORDER:
			ActionStartOrder(req, resp)
    case ACTION_CHECK_TIME_CLOSE:
      bot.ActionCheckTimeClose(req, resp)
		/*case ACTION_START_ORDER_WITH_INPUTS:
			bot.ActionStartOrderWithInputs(reqObj, rasaResp)
    case ACTION_UPDATE_ORDER:
      bot.ActionUpdateOrder(reqObj, rasaResp)
    case ACTION_CHECK_IS_OPEN:
      bot.actionUtter(reqObj, "Yes")
    case ACTION_CHECK_IS_OPEN_ON_DAY:
      bot.actionUtter(reqObj, "Yes")
    case ACTION_CHECK_TIME_CLOSE:
      bot.actionUtter(reqObj, "Yes")
    case ACTION_CHECK_TIME_CLOSE_ON_DAY:
      bot.actionUtter(reqObj, "Yes")
    case ACTION_RESET_SLOTS:
      bot.ActionResetSlots(reqObj)*/
		}

    return resp, nil
}


func ActionStartOrder(req *RasaRequest, resp *RasaResponse) {
	order := Order{
		RecipientId:          req.SenderId,
    RecipientContact: req.Tracker.Slots["recipient_contact"],
		StartTime:            currentTime(),
		LastModificationTime: currentTime(),
		IsVisible:            true,
	}

	slots := req.Tracker.Slots
	if slots["address"] != "" {
		order.Address = slots["address"]
	}

	if slots["name"] != "" {
		order.Name = slots["name"]
	}

	if slots["type"] != "" {
		order.Type = slots["type"]
	}

	if slots["content"] != "" {
		order.Content = slots["content"]
	}

	bot.saveOrder(req, &order)
}

func (bot *Bot) ActionUpdateOrder(req *RasaRequest, resp *RasaResponse) {
  businessId := req.Tracker.Slots["business_id"]
  recipientId := req.Tracker.Slots["recipient_id"]

	orderQuery := bot.Client.Collection(Businesses).Doc(businessId).Collection(Orders).Where("recipientId", "==", recipientId)
  orderQuery = orderQuery.Where("visible", "==", true).OrderBy("lastModificationTime", firestore.Desc)

  order := &Order{}
	iter := orderQuery.Documents(bot.Ctx)
	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Println(err)
			// TODO handle error
			break
		}

		err = doc.DataTo(order)
		if err != nil {
			log.Println(err)
		}

		order.Id = doc.Ref.ID
  }

  if order.Id == "" {
    // TODO No order, so will make a new one
    reply := fmt.Sprintf("")
    resp.Responses = append(resp.Responses, Response{Text: reply})
  } else {
    orderRef := bot.Client.Collection(Businesses).Doc(businessId).Collection(Orders).Doc(order.Id)

    // Now we have our order and can update it
    orderRef.Update(bot.Ctx, []firestore.Update{
      {Path: "address", Value: req.Tracker.Slots["address"]},
      {Path: "name", Value: req.Tracker.Slots["name"]},
      {Path: "type", Value: req.Tracker.Slots["type"]},
      {Path: "content", Value: req.Tracker.Slots["content"]},
      {Path: "lastModifiedTime", Value: currentTime()},
    })

  }


}

func (bot *Bot) ActionCheckTimeClose(req *RasaRequest, resp *RasaResponse) {
  // Need to check the database to see if business is closed or not
  // Then modifies the RasaResponse with the correct RasaResponse

  businessId := req.Tracker.Slots["business_id"]
  dataSnap, err := bot.Client.Collection(Businesses).Doc(businessId).Get(bot.Ctx)

  if err != nil {
    log.Println(err)
  }

  var business Business
  err = dataSnap.DataTo(&business)

  if err != nil {
    log.Println(err)
  }

  // TODO make the response dynamic
  reply := fmt.Sprintf("We close at %s", business.TimeClose())
  resp.Responses = append(resp.Responses, Response{Text: reply})
}

func (bot *Bot) ActionCheckTimeCloseOnDay(req *RasaRequest, resp *RasaResponse) {
  // Need to check the database to see if business is closed or not
  // Then modifies the RasaResponse with the correct RasaResponse
  businessId := req.Tracker.Slots["business_id"]
  dataSnap, err := bot.Client.Collection(Businesses).Doc(businessId).Get(bot.Ctx)

  if err != nil {
    log.Println(err)
  }

  var business Business
  err = dataSnap.DataTo(&business)

  if err != nil {
    log.Println(err)
  }

  var entity Entity
  for _, v := range req.Tracker.LatestMessage.Entities {
    if v.Entity == "time" {
      entity = v
      break
    }
  }

  t, err := time.Parse(time.RFC3339, entity.Value)

  if err != nil {
    log.Println(err)
  }

  // TODO make the response dynamic
  reply := fmt.Sprintf("On %d/%d we close at %s", t.Month(), t.Day(), business.TimeClose())
  resp.Responses = append(resp.Responses, Response{Text: reply})
}

func (bot Bot) saveOrder(req *RasaRequest, order *Order) {
  businessId := req.Tracker.Slots["business_id"]
  recipientId := req.Tracker.Slots["recipient_id"]


	ordersRef := bot.Client.Collection(Businesses).Doc(businessId).Collection(Orders)

	docRef, _, err := ordersRef.Add(bot.Ctx, order)

	if err != nil {
		log.Println(err)
	}

	order.Id = docRef.ID

	recipientRef := bot.Client.Collection(Businesses).Doc(businessId).Collection(Recipients).Doc(recipientId)
	recipientRef.Update(ctx, []firestore.Update{
		{Path: RecentOrderId, Value: order.Id},
	})

}

func currentTime() int64 {
	return time.Now().UnixNano() / 1000000
}

func (bot *Bot) saveMessage(business *Business, recipeint *Recipient, message *Message) error {
	messagesRef := bot.Client.Collection(Businesses).Doc(business.Id).Collection(Messages)
	docRef, _, err := messagesRef.Add(bot.Ctx, message)

	if err != nil {
		return err
	}

	message.Id = docRef.ID
	return nil
}

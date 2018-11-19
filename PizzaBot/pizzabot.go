package main

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"math"
	"net/http"
	"time"
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

	// Need to save the new message to firebase
	err := bot.saveMessage(reqObj.Business, reqObj.Recipient, reqObj.Message)

	if err != nil {
		log.Println(err)
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
	log.Println(req.Tracker.Slots)
	switch action {
	case ACTION_UPDATE_ORDER:
		bot.ActionUpdateOrder(req, resp)
	case ACTION_CHECK_IS_OPEN:
		bot.ActionCheckIsOpen(req, resp)
	case ACTION_CHECK_IS_OPEN_ON_DAY:
		bot.ActionCheckIsOpenOnDay(req, resp)
	case ACTION_CHECK_TIME_CLOSE:
		bot.ActionCheckTimeClose(req, resp)
	case ACTION_CHECK_TIME_CLOSE_ON_DAY:
		bot.ActionCheckTimeCloseOnDay(req, resp)
	case ACTION_CHECK_RESERVATION_DATETIME:
		bot.ActionCheckReservationDatetime(req, resp)
	case ACTION_SET_SCHEDULED_TIME_SLOT:
		bot.ActionSetScheduledTimeSlot(req, resp)
	case ACTION_SET_SIZE_SLOT:
		bot.ActionSetSizeSlot(req, resp)
	}

	return resp, nil
}

func (bot *Bot) ActionCheckReservationDatetime(req *RasaRequest, resp *RasaResponse) {
	// TODO add support for time intervals
	/* Will have a datetime, businessId, partySize, etc saved in slots */
	businessId := req.Tracker.Slots["business_id"]
	//recipientId := req.Tracker.Slots["recipient_id"]
	searchTimeStr := req.Tracker.Slots["time"]
	partySize := req.Tracker.Slots["partySize"]
	name := req.Tracker.Slots["name"]

	searchTime, err := time.Parse(time.RFC3339, searchTimeStr)

	if err != nil {
		log.Println(err)
	}

	business, err := bot.getBusinessFromId(businessId)

	if err != nil {
		log.Println(err)
	}

	reservationResult, err := Query(business.ReservationPlatformId, searchTime, partySize)

	if err != nil {
		log.Println(err)
	}

	if reservationResult.Message == "" {
		// Reservations found within 2.5 hours of request

		found := false
		// Check if one equals exactly and if so make the reservation
		for _, v := range reservationResult.Results {
			if v == searchTime {
				found = true
			}
		}

		if found {
			// Exact match, so as long as we have a name we add the reservation to the db

			if name == "" {
				// Action ask the name
				nextAction := Event{Event: "followup", Name: "utter_ask_name"}
				resp.Events = append(resp.Events, nextAction)
				return
			} else {
				// Force Action save_reservation
				nextAction := Event{Event: "followup", Name: "action_save_reservation"}
				resp.Events = append(resp.Events, nextAction)
				return
			}
		}

		// Check if any are within 15 minutes and if so ask if that is fine
		lessThan15 := false
		selectedTime := reservationResult.Results[0]
		for _, v := range reservationResult.Results {
			if math.Abs(v.Sub(searchTime).Minutes()) <= 15 {
				lessThan15 = true
				selectedTime = v
			}
		}

		if lessThan15 {
			// Action is this one good?
			nextAction := Event{Event: "slot", Name: "potential_times", Value: []time.Time{selectedTime}}
			resp.Events = append(resp.Events, nextAction)

			nextAction = Event{Event: "followup", Name: "action_ask_is_close_time_okay"}
			resp.Events = append(resp.Events, nextAction)
			return
		}

		nextAction := Event{Event: "slot", Name: "potential_times", Value: reservationResult.Results}
		resp.Events = append(resp.Events, nextAction)

		// action we didn't find any at that time, but do any of these times work for you?
		nextAction = Event{Event: "followup", Name: "action_ask_if_any_similar_times_work"}
		resp.Events = append(resp.Events, nextAction)
		return

	} else if reservationResult.Message == NO_AVAILABLE {
		nextAction := Event{Event: "followup", Name: "utter_no_reservations_available"}
		resp.Events = append(resp.Events, nextAction)
		return

	} else if reservationResult.Message == IN_ADVANCE {
		nextAction := Event{Event: "followup", Name: "utter_time_within"}
		resp.Events = append(resp.Events, nextAction)
		return
	}

}

func (bot *Bot) ActionSetScheduledTimeSlot(req *RasaRequest, resp *RasaResponse) {
	scheduledTime := ""
	for _, v := range req.Tracker.LatestMessage.Entities {
		if v.Entity == "time" {
			scheduledTime = v.Value.(string)
		}
	}
	nextAction := Event{Event: "slot", Name: "scheduledTime", Value: scheduledTime}
	resp.Events = append(resp.Events, nextAction)
}

func (bot *Bot) ActionSetSizeSlot(req *RasaRequest, resp *RasaResponse) {
	size := ""
	for _, v := range req.Tracker.LatestMessage.Entities {
		if v.Entity == "number" {
			size = v.Value.(string)
		}
	}
	nextAction := Event{Event: "slot", Name: "size", Value: size}
	resp.Events = append(resp.Events, nextAction)
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

	orderType := req.Tracker.Slots["type"]

	if req.Tracker.Slots["address"] != "" && orderType == "" {
		orderType = "DELIVERY"
		resp.Events = append(resp.Events, Event{"slot", "type", "DELIVERY"})
	}

	if order.Id == "" {
		order := Order{
			RecipientId:          req.SenderId,
			RecipientContact:     req.Tracker.Slots["recipient_contact"],
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

		order.Type = orderType

		if slots["contents"] != "" {
			order.Content = slots["contents"]
		}
		bot.saveOrder(req, &order)

	} else {
		orderRef := bot.Client.Collection(Businesses).Doc(businessId).Collection(Orders).Doc(order.Id)

		// Now we have our order and can update it
		orderRef.Update(bot.Ctx, []firestore.Update{
			{Path: "address", Value: req.Tracker.Slots["address"]},
			{Path: "name", Value: req.Tracker.Slots["name"]},
			{Path: "type", Value: orderType},
			{Path: "content", Value: req.Tracker.Slots["contents"]},
			{Path: "lastModificationTime", Value: currentTime()},
		})

	}

}

func (bot *Bot) ActionAskNext(req *RasaRequest, resp *RasaResponse) {
	// Figure out which question to ask and return it as a follow up action

	slots := req.Tracker.Slots
	emptySlots := map[string]string{}

	for k, v := range slots {
		if v == "" {
			emptySlots[k] = v
		}
	}

	// Now we have our empty slots and we can probably just choose any
}

func (bot *Bot) ActionCheckTimeClose(req *RasaRequest, resp *RasaResponse) {
	// Need to check the database to see if business is closed or not
	// Then modifies the RasaResponse with the correct RasaResponse

	businessId := req.Tracker.Slots["business_id"]
	business, err := bot.getBusinessFromId(businessId)

	if err != nil {
		log.Println(err)
	}

	// TODO make the response dynamic
	reply := fmt.Sprintf("We close at %s", business.TimeClose(""))
	resp.Responses = append(resp.Responses, Response{Text: reply})
}

func (bot *Bot) ActionCheckTimeCloseOnDay(req *RasaRequest, resp *RasaResponse) {
	// Need to check the database to see if business is closed or not
	// Then modifies the RasaResponse with the correct RasaResponse
	businessId := req.Tracker.Slots["business_id"]
	business, err := bot.getBusinessFromId(businessId)

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

	t, err := time.Parse(time.RFC3339, entity.Value.(string))

	if err != nil {
		log.Println(err)
	}

	// TODO make the response dynamic
	reply := fmt.Sprintf("On %d/%d we close at %s", t.Month(), t.Day(), business.TimeClose(""))
	resp.Responses = append(resp.Responses, Response{Text: reply})
}

func (bot *Bot) ActionCheckIsOpen(req *RasaRequest, resp *RasaResponse) {
	businessId := req.Tracker.Slots["business_id"]
	business, err := bot.getBusinessFromId(businessId)

	if err != nil {
		log.Println(err)
	}

	// TODO make dynamic
	reply := ""
	if business.IsOpen() {
		reply = "Yes, we're open"
	} else {
		reply = "Sorry, no we're not"
	}
	resp.Responses = append(resp.Responses, Response{Text: reply})
}

func (bot *Bot) ActionCheckIsOpenOnDay(req *RasaRequest, resp *RasaResponse) {
	businessId := req.Tracker.Slots["business_id"]
	business, err := bot.getBusinessFromId(businessId)

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

	t, err := time.Parse(time.RFC3339, entity.Value.(string))

	if err != nil {
		log.Println(err)
	}

	// TODO make dynamic
	reply := ""
	if business.IsOpen() {
		reply = "Yes, we're open"
	} else {
		reply = "Sorry, no we're not"
	}

	reply = fmt.Sprint("%s on %d/%d", reply, t.Month(), t.Day())
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

func (bot *Bot) getBusinessFromId(businessId string) (Business, error) {
	var business Business
	dataSnap, err := bot.Client.Collection(Businesses).Doc(businessId).Get(bot.Ctx)

	if err != nil {
		return business, err
	}

	err = dataSnap.DataTo(&business)

	if err != nil {
		return business, err
	}

	return business, nil

}

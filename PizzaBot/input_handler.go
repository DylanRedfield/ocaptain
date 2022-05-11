package main

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"log"
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

	smsRequest := MessageRequest{
		To:   reqObj.Recipient.Contact,
		From: reqObj.Business.PhoneNumber,
		Body: reqObj.Message}

	log.Println(reqObj.Business)

	if reqObj.Recipient.Platform == FACEBOOK_MESSENGER_PLATFORM {
		smsRequest.From = reqObj.Business.FacebookMessengerId
		reqObj.Business.FacebookMessengerClient.Send(&smsRequest)
	} else if reqObj.Recipient.Platform == TWILIO_WHATSAPP_PLATFORM {
		log.Println(reqObj.Business.Whatsapp)
		smsRequest.From = reqObj.Business.Whatsapp
		reqObj.Business.TwilioClient.Send(&smsRequest)
	} else if reqObj.Business.SmsPlatform == "TWILIO" {
		reqObj.Business.TwilioClient.Send(&smsRequest)
	} else if reqObj.Business.SmsPlatform == "SWIFT" {
		log.Println("Swift send")
		bot.SwiftClient.Send(&smsRequest)
	}

	return BusinessResponse{}
}

func (bot *Bot) HandleOutsideInput(reqObj *OutsideRequest) OutsideResponse {

	businessId := reqObj.Business.Id

	// Need to check if a recipient was found, and if not create one, and if so update the recent message
	if reqObj.Recipient.Id == "" {
		reqObj.Recipient.RecentMessage = reqObj.Message

		personRef, _, err := bot.Client.Collection(Businesses).Doc(businessId).Collection(Recipients).Add(bot.Ctx, reqObj.Recipient)
		reqObj.Recipient.Id = personRef.ID
		reqObj.Message.RecipientId = personRef.ID

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

	//bot.sendToAI(reqObj)

	if bot.IsDemo {
		bot.sanderDemo(reqObj)
	}
	if reqObj.Business.SmsNotifyEnabled {
		log.Println("Noitify Enabled")
		bot.notifyStaff(reqObj)
	}
	log.Print("end")
	log.Println(time.Now().String())
	return OutsideResponse{}
}

func (bot *Bot) notifyStaff(reqObj *OutsideRequest) {
	var employees = reqObj.Business.Employees
	actives := []string{}

	for _, employee := range employees {
		log.Println("Active Employee")
		if employee.IsActive {
			log.Println("Active Employee")
			actives = append(actives, employee.PhoneNumber)
		}
	}

	bulkReq := &BulkMessageRequest{actives, reqObj.Message.Content}
	reqObj.Business.TwilioClient.SendBulk(bulkReq)
}

func (bot *Bot) sanderDemo(receivedMsg *OutsideRequest) {
	responses := []string{
		"Sure, for what time?",
		"Can you do 20:30?",
		"And what's you name?",
		"Great, see you then",
		"No, we serve meat and dairy",
	}

	// Send message to customer via whatsapp client

	generatedBody := responses[bot.DemoCounter]
	bot.DemoCounter = (bot.DemoCounter + 1) % len(responses)

	msgReq := MessageRequest{
		To:   receivedMsg.Recipient.Contact,
		Body: generatedBody,
	}
	msg := Message{}
	msg.Content = generatedBody
	msg.DidBotCreate = true
	msg.HasBusinessRead = false
	msg.RecipientId = receivedMsg.Recipient.Id
	msg.IsBusinessSender = true
	msg.TimeSent = time.Now().UnixNano() / 1000000
	if receivedMsg.Recipient.Platform == TWILIO_WHATSAPP_PLATFORM {
		msgReq.Platform = TWILIO_WHATSAPP_PLATFORM
		msgReq.From = receivedMsg.Business.Whatsapp
	} else {
		msgReq.From = receivedMsg.Business.PhoneNumber
	}

	receivedMsg.Business.TwilioClient.Send(&msgReq)

	bot.saveMessage(receivedMsg.Business, receivedMsg.Recipient, &msg)
}

func (bot *Bot) sendToAI(reqObj *OutsideRequest) OutsideResponse {
	// Send a http request that will be handled in the textual_input_channel
	// The body is the OutsideRequest object
	body, err := json.Marshal(reqObj)

	if err != nil {
		log.Println(err)
	}

	//envValues := GetEnvValues()

	rasaUrl := fmt.Sprintf("http://localhost:%s/webhooks/textual/webhook", "5005")
	req, err := http.NewRequest("POST", rasaUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Println(err)
	}

	http.DefaultClient.Do(req)

	return OutsideResponse{}

}

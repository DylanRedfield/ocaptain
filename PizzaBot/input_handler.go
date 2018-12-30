package main

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
  "os"
  "io/ioutil"
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

	jsonFile, err := os.Open("../env_values.json")

  if err != nil {
    log.Println(err)
  }

  defer jsonFile.Close()

  byteValue, _:= ioutil.ReadAll(jsonFile)

  var envValues EnvValues

	json.Unmarshal([]byte(byteValue), &envValues)
	rasaUrl := fmt.Sprintf("http://localhost:%s/webhooks/textual/webhook", envValues.RasaPort)
	req, err := http.NewRequest("POST", rasaUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Println(err)
	}

	http.DefaultClient.Do(req)

	return OutsideResponse{}
}

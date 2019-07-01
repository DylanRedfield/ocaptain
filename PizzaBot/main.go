package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var ctx context.Context

var bot *Bot

func init() {
	ctx = context.Background()

	var err error
	bot, err = NewBot(ctx)

	if err != nil {
		log.Println("Error initiating bot")
	}
}

func main() {

	name := GetEnvValues().Name

	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("/%s/businessInput", name), http.HandlerFunc(businessInput))
	mux.Handle(fmt.Sprintf("/%s/outsideSmsInput", name), http.HandlerFunc(outsideSmsInput))

	jsonFile, err := os.Open("../env_values.json")

	if err != nil {
		log.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var envValues EnvValues

	json.Unmarshal([]byte(byteValue), &envValues)

	log.Println(envValues.PizzaPort)
	log.Println(http.ListenAndServe(":"+envValues.PizzaPort, mux))
}

// Recieves a BotRequest as HTTP payload,
// and returns BotResponse as HTTP payload.
func businessInput(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		// TODO handle error
	}

	var reqObj BusinessRequest

	if err := json.Unmarshal(body, &reqObj); err != nil {
		// TODO handle error
	}

	bot.HandleBusinessInput(reqObj)
}

// Recieves input from SMS service like Twilio
func outsideSmsInput(w http.ResponseWriter, req *http.Request) {
	// Params come in from GET URL.
	// I can get them as a map, but not obj.
	// So I marshal the map into a json string,
	// then unmarshal the json shring into the object

	values := req.URL.Query()
	reqObj := MessageRequest{}
	log.Println(values)
	if val, exists := values["smsPlatform"]; exists {
		if val[0] == "SWIFT" {
			to := fmt.Sprintf("+1%s", values["Destination"][0])
			from := fmt.Sprintf("+%s", values["PhoneNumber"][0])
			reqObj = MessageRequest{To: to, Body: values["MessageBody"][0], From: from}
		}
	} else {
		reqObj = MessageRequest{To: req.URL.Query()["To"][0], Body: req.URL.Query()["Body"][0], From: req.URL.Query()["From"][0]}
	}
	// TODO will error on swift message from conflicting names

	outsideReq := toOutsideRequest(reqObj)
	bot.HandleOutsideInput(&outsideReq)
}



// Turns TwilioRequest into standard OutsideRequest object
func toOutsideRequest(twilReq MessageRequest) OutsideRequest {

	timeInMil := time.Now().UnixNano() / 1000000
	message := &Message{
		Content:          twilReq.Body,
		IsBusinessSender: false,
		HasBusinessRead:  false,
		DidBotCreate:     false,
		TimeSent:         timeInMil,
	}

	business, err := businessFromPhone(twilReq.To)

	if err != nil {
		log.Println(err)
	}

	recipient, err := recipientFromNumber(twilReq.From, business.Id)

	if err != nil {
		log.Println(err)
	}

	message.RecipientId = recipient.Id

	return OutsideRequest{Recipient: recipient, Message: message, Business: business}
}

func businessFromPhone(phoneNumber string) (*Business, error) {
	business := &Business{}
	log.Println(phoneNumber)
	iter := bot.Client.Collection(Businesses).Where(PhoneNumber, "==", phoneNumber).Documents(bot.Ctx)

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

		err = doc.DataTo(business)

		if err != nil {
			log.Println(err)
		}

		business.Id = doc.Ref.ID

		// Doesnt automatically unmarshall subcollections
		subIter := bot.Client.Collection(Businesses).Doc(business.Id).Collection("employees").Documents(bot.Ctx)
		employees := []Employee{}
		for {
			subDoc, err := subIter.Next()

			if err == iterator.Done {
				break
			}
			if err != nil {
				return business, err
			}

			employee := Employee{}
			err = subDoc.DataTo(&employee)

			if err != nil {
				return business, err
			}

			employee.Id = subDoc.Ref.ID

			employees = append(employees, employee)
		}
		business.Employees = employees
		log.Println(employees)

	}

	if business.Id == "" {
		return business, errors.New("Business not found")
	} else {
		return business, nil
	}

}

// Takes in recipientId and returns recipient or error if none is found or there was an error retrieving data
func recipientFromNumber(recipientNumber string, businessId string) (*Recipient, error) {

	query := bot.Client.Collection(Businesses).Doc(businessId).Collection(Recipients).Where(Contact, "==", recipientNumber)

	iter := query.Documents(bot.Ctx)

	recipient := &Recipient{}

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

		err = doc.DataTo(recipient)

		if err != nil {
			log.Println(err)
			break
		}

		recipient.Id = doc.Ref.ID
	}

	// Check if recipient was found, because ID (or any value) will be empty
	if recipient.Id == "" {
		recipient.Contact = recipientNumber
		// No recipient found, so return error
		return recipient, errors.New("No matching recipient founder")
	} else {
		recipient.Contact = recipientNumber
		return recipient, nil
	}
}

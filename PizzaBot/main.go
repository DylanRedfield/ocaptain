package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
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

	mux := http.NewServeMux()
	mux.Handle("/PizzaBot/businessInput", http.HandlerFunc(businessInput))
	mux.Handle("/PizzaBot/outsideSmsInput", http.HandlerFunc(outsideSmsInput))
	mux.Handle("/PizzaBot/sendSelf", http.HandlerFunc(sendSelf))
	mux.Handle("/ocaptain", http.HandlerFunc(actionInput))
	mux.Handle("/ocaptain/sendAndSave", http.HandlerFunc(sendAndSave))

	jsonFile, err := os.Open("../env_values.json")

	if err != nil {
		log.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var envValues EnvValues

	json.Unmarshal([]byte(byteValue), &envValues)

	log.Println(http.ListenAndServe(":"+envValues.PizzaPort, mux))
}

func test() {
	datetime := time.Date(2018, 11, 19, 14, 0, 0, 0, time.Local)
	result, err := Query("24712", datetime, "3")

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%s", result.Results)
}

func actionInput(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Println(err)
	}

	var reqObj RasaRequest
	if err := json.Unmarshal(body, &reqObj); err != nil {
		log.Println(err)
	}

	resp, err := bot.HandleAction(&reqObj)

	if err != nil {
		log.Print(err)
	}

	respString, err := json.Marshal(*resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(respString)

	if err != nil {
		log.Println(err)
	}

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

  log.Println("Recieved")
	// TODO will error on swift message from conflicting names
	reqObj := MessageRequest{To: req.URL.Query()["To"][0], Body: req.URL.Query()["Body"][0], From: req.URL.Query()["From"][0]}

	outsideReq := toOutsideRequest(reqObj)
	bot.HandleOutsideInput(&outsideReq)
}

func sendSelf(w http.ResponseWriter, req *http.Request) {
	reqObj := MessageRequest{To: "+12027593168", Body: "Default message", From: "+12027593168"}
	outsideReq := toOutsideRequest(reqObj)
	bot.HandleOutsideInput(&outsideReq)
}

func sendAndSave(w http.ResponseWriter, req *http.Request) {
	log.Println("UGHHH")
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Println(err)
	}

	var reqObj OutsideRequest
	if err := json.Unmarshal(body, &reqObj); err != nil {
		log.Println(err)
	}

	twilioClient := TwilioClient{
		AccountSid: "AC9dfbda388f3ee10353bbc001694f5c27",
		AuthToken:  "e3429e06cc27740f1c859d2bfc9964ae"}

	to := reqObj.Recipient.Contact
	from := reqObj.Business.PhoneNumber
	text := reqObj.Message.Content

	twilioClient.Send(&MessageRequest{to, from, text})

	ctx = context.Background()

	sa := option.WithCredentialsFile("firebase-config.json")

	app, err := firebase.NewApp(ctx, nil, sa)

	if err != nil {
		log.Println(err)
	}

	client, err := app.Firestore(ctx)

	if err != nil {
		log.Println(err)
	}

	log.Println(reqObj.Business.Id)
	messagesRef := client.Collection(Businesses).Doc(reqObj.Business.Id).Collection(Messages)
	_, _, err = messagesRef.Add(ctx, reqObj.Message)

	if err != nil {
		log.Println(err)
	}

	personRef := client.Collection(Businesses).Doc(reqObj.Business.Id).Collection(Recipients).Doc(reqObj.Recipient.Id)
	personRef.Update(ctx, []firestore.Update{
		{Path: RecentMessage, Value: reqObj.Message},
	})

}
func initFirebase() *Bot {
	ctx = context.Background()
	bot, err := NewBot(ctx)

	if err != nil {
		// TODO handle error
	}

	return bot
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

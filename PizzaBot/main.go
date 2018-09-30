package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var ctx context.Context

func main() {
	mux := http.NewServeMux()
	mux.Handle("/PizzaBot/businessInput", http.HandlerFunc(businessInput))
	mux.Handle("/PizzaBot/outsideSmsInput", http.HandlerFunc(outsideSmsInput))
	mux.Handle("/PizzaBot/sendSelf", http.HandlerFunc(sendSelf))
  mux.Handle("/Textual/startOrder", http.HandlerFunc(ActionStartOrder))
	log.Println(http.ListenAndServe(":8080", mux))
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

	HandleBusinessInput(ctx, reqObj)
}

// Recieves input from SMS service like Twilio
func outsideSmsInput(w http.ResponseWriter, req *http.Request) {
	// Params come in from GET URL.
	// I can get them as a map, but not obj.
	// So I marshal the map into a json string,
	// then unmarshal the json shring into the object

	reqObj := TwilioRequest{To: req.URL.Query()["To"][0], Body: req.URL.Query()["Body"][0], From: req.URL.Query()["From"][0]}

	outsideReq := ToOutsideRequest(reqObj)
  log.Println("Phone", outsideReq.Recipient.Contact)
	HandleOutsideInput(ctx, outsideReq)
}

func sendSelf(w http.ResponseWriter, req *http.Request) {
	reqObj := TwilioRequest{To: "+12027593168", Body: "Default message", From: "+12027593168"}
	outsideReq := ToOutsideRequest(reqObj)
	HandleOutsideInput(ctx, outsideReq)
}

func initFirebase() Bot {
	ctx = context.Background()
	bot, err := NewBot(ctx)

	if err != nil {
		// TODO handle error
	}

	return bot
}

// Turns TwilioRequest into standard OutsideRequest object
func ToOutsideRequest(twilReq TwilioRequest) OutsideRequest {

	timeInMil := time.Now().UnixNano() / 1000000
	message := &Message{
    Content: twilReq.Body, 
    IsBusinessSender: false,
		HasBusinessRead: false, 
    DidBotCreate: false, 
    TimeSent: timeInMil,
  }

	//recipient := &Recipient{Contact: twilReq.From, Platform: "SMS"}

	sa := option.WithCredentialsFile("firebase-config.json")

	app, err := firebase.NewApp(ctx, nil, sa)

	if err != nil {
		log.Println(err)
	}

	client, err := app.Firestore(ctx)

	if err != nil {
		log.Println(err)
	}

	business, err := businessFromPhone(client, ctx, twilReq.To)

	if err != nil {
		log.Println(err)
	}

	recipient, err := recipientFromNumber(client, ctx, message, twilReq.From, business.Id)

	if err != nil {
		log.Println(err)
	}

	message.RecipientId = recipient.Id

	err = saveMessage(client, ctx, business, recipient, message)

	if err != nil {
		log.Println(err)
	}

	return OutsideRequest{Recipient: recipient, Message: message, Business: business}
}

func saveMessage(client *firestore.Client, ctx context.Context, business *Business, recipeint *Recipient,
	message *Message) error {
	messagesRef := client.Collection(Businesses).Doc(business.Id).Collection(Messages)
	docRef, _, err := messagesRef.Add(ctx, message)

	if err != nil {
		return err
	}

	message.Id = docRef.ID
	return nil
}

func businessFromPhone(client *firestore.Client, ctx context.Context, phoneNumber string) (*Business, error) {
	business := &Business{}
	iter := client.Collection(Businesses).Where(PhoneNumber, "==", phoneNumber).Documents(ctx)

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
		business.Id = doc.Ref.ID

		if err != nil {
			log.Println(err)
		}
	}

	if business.Id == "" {
		return business, errors.New("Business not found")
	} else {
		return business, nil
	}

}

// Takes in recipientId
func recipientFromNumber(client *firestore.Client, ctx context.Context, message *Message, recipientNumber string, businessId string) (*Recipient, error) {

	query := client.Collection(Businesses).Doc(businessId).Collection(Recipients).Where(Contact, "==", recipientNumber)

	iter := query.Documents(ctx)

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
		}

		recipient.Id = doc.Ref.ID
	}

	if recipient.Id == "" {
		recipient.RecentMessage = message
    recipient.Contact = recipientNumber

		// No recipient was found in firebase, so need to construct a new one
		personRef, _, err := client.Collection(Businesses).Doc(businessId).Collection(Recipients).Add(ctx, recipient)
		recipient.Id = personRef.ID
		return recipient, err
	} else {
		personRef := client.Collection(Businesses).Doc(businessId).Collection(Recipients).Doc(recipient.Id)
		personRef.Update(ctx, []firestore.Update{
			{Path: RecentMessage, Value: message},
		})
		return recipient, nil
	}

}

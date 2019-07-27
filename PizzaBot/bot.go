package main

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"time"
)

type Bot struct {
	Client       *firestore.Client
	Ctx          context.Context
	TwilioClient TwilioClient
	SwiftClient  SwiftClient
	ActiveMessages []*OutsideRequest
}

func NewBot(ctx context.Context) (*Bot, error) {
	sa := option.WithCredentialsFile("firebase-config.json")

	app, err := firebase.NewApp(ctx, nil, sa)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	client, err := app.Firestore(ctx)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	envValues := GetEnvValues()
	twilioClient := TwilioClient{
		AccountSid: envValues.TwilioAccountSid,
		AuthToken:  envValues.TwilioAuthToken,
	}

	swiftClient := SwiftClient{
		AccountKey: envValues.SwiftAccountKey}

	return &Bot{Client: client, Ctx: ctx, TwilioClient: twilioClient, SwiftClient: swiftClient}, nil
}

func (bot *Bot) CheckActiveMessages() {

	for _, request := range bot.ActiveMessages {
		message := request.Message
		timeSent := time.Unix(message.TimeSent, 0)

		since := time.Since(timeSent)

		if int64(since) / 60 >= 5 {
			// Time has ellapsed so send message to the recipient from the correct number
			// and remove the message from the list

			text := fmt.Sprintf("Sorry we're not about to respond right now. Could you give us a call" +
				"at %s", request.Business.PhoneNumber)
			messageReq := MessageRequest{To: request.Recipient.Contact, Body: text, From: request.Business.PhoneNumber}
			if request.Business.SmsPlatform == "SWIFT" {
				bot.SwiftClient.Send(&messageReq)
			} else if request.Business.SmsPlatform == "TWILIO" {
				bot.TwilioClient.Send(&messageReq)
			}
		}

	}

}

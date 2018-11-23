package main

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
)

type Bot struct {
	Client    *firestore.Client
	Ctx       context.Context
	SmsClient TwilioClient
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

	twilioClient := TwilioClient{
		AccountSid: "AC9dfbda388f3ee10353bbc001694f5c27",
		AuthToken:  "e3429e06cc27740f1c859d2bfc9964ae"}

	return &Bot{Client: client, Ctx: ctx, SmsClient: twilioClient}, nil
}

type BusinessRequest struct {
	BusinessId string
	Business   Business
	Message    string
	Recipient  Recipient
	IsPhone    bool
}

type BusinessResponse struct {
	BusinessId  string
	Message     string
	RecipientId string
}

type Recipient struct {
	Id            string   `firestore:"-"`
	Name          string   `firestore:"name,omitempty"`
	Address       string   `firestore:"name,omitempty"`
	Contact       string   `firestore:"contact"`
	Platform      Platform `firestore:"platform,omitempty"`
	RecentMessage *Message `firestore:"recentMessage,omitempty"`
	RecentOrderId string   `firestore:"recentOrderId,omitempty"`
}

type Platform string

type Message struct {
	Id               string `firestore:"-"`
	Content          string `firestore:"content"`
	IsBusinessSender bool   `firestore:"isBusinessSender"`
	TimeSent         int64  `firestore:"timeSent"`
	DidBotCreate     bool   `firestore:"didBotCreate"`
	HasBusinessRead  bool   `firestore:"hasBusinessRead"`
	RecipientId      string `firestore:"recipientId"`
}

type OutsideRequest struct {
	Id        string     `json:"id"`
	Recipient *Recipient `json:"recipient"`
	Message   *Message   `json:"message"`
	Business  *Business  `json:"business"`
}

type OutsideResponse struct {
}

type Order struct {
	Id                   string `firestore:"-"`
	Address              string `firestore:"address"`
	Name                 string `firestore:"name"`
	Content              string `firestore:"content"`
	StartTime            int64  `firestore:"startTime"`
	CompleteTime         int64  `firestore:"completeTime"`
	ScheduledTime        int64  `firestore:"scheduledTime"`
	LastModificationTime int64  `firestore:"lastModificationTime"`
	Type                 string `firestore:"type"`
	IsVisible            bool   `firestore:"visible"`
	RecipientId          string `firestore:"recipientId"`
	RecipientContact     string `firestore:"recipientContact"`
}

type Business struct {
	Id                    string               `firestore:"-"`
	Approved              bool                 `firestore:"approved"`
	Password              string               `firestore:"password"`
	PhoneNumber           string               `firestore:"phoneNumber"`
	Hours                 map[string]OpenClose `firestore:"hours"`
	ReservationPlatform   string               `firestore:"reservationPlatform"`
	ReservationPlatformId string               `firestore:"reservationPlatformId"`
}

type OpenClose struct {
	IsOpen    bool  `firestore:"isOpen"`
	OpenTime  int32 `firestore:"openTime"`
	CloseTime int32 `firestore:"closeTime"`
}

type Reservation struct {
	Id            string `firestore:"-"`
	RecipientId   string `firestore:"recipientId"`
	Name          string `firestore:"-"`
	ScheduledTime int64  `firestore:"scheduledTime"`
	NumPeople     string  `firestore:"numPeople"`
	IsVisible     bool   `firestore:"isVisible"`
}

func (business Business) TimeClose(day string) int32 {
	return business.Hours[day].CloseTime
}

func (business Business) IsOpen() bool {
	// TODO implement
	return true
}

type Tracker struct {
	Slots         map[string]interface{} `json:"slots"`
	SenderId      string                 `json:"sender_id"`
	LatestMessage LatestMessage          `json:"latest_message"`
}

type LatestMessage struct {
	Text string `json:"text"`
	//	Intent   string   `json:"intent"`
	Entities []Entity `json:"entities"`
}

type Entity struct {
	Start      int32       `json:"start"`
	End        int32       `json:"end"`
	Value      interface{} `json:"value"`
	Text       string      `json:"text"`
	Confidence float64     `json:"confidence"`
	Entity     string      `json:"entity"`
}

type RasaRequest struct {
	NextAction string  `json:"next_action"`
	SenderId   string  `json:"sender_id"`
	Tracker    Tracker `json:"tracker"`
}

type RasaResponse struct {
	Events    []Event    `json:"events,omitempty"`
	Responses []Response `json:"responses"`
}

func NewRasaResponse() *RasaResponse {
	return &RasaResponse{
		Events:    []Event{},
		Responses: []Response{},
	}
}

type Response struct {
	Text string `json:"text"`
}

type Event struct {
	Event string      `json:"event"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

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

type RasaTime struct {
	Value string
	Grain string
	Type  string
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
	Hours                 map[int]OpenClose    `firestore:"hours"`
	HoursExcpetions       map[string]OpenClose `firestore:"hoursExceptions"`
	ReservationPlatform   string               `firestore:"reservationPlatform"`
	ReservationPlatformId string               `firestore:"reservationPlatformId"`
}

type OpenClose struct {
	IsOpen bool  `firestore:"isOpen"`
	Open   int64 `firestore:"open"`
	Close  int64 `firestore:"close"`
}

func (openClose OpenClose) ClosePastMidnight() bool {
	return openClose.Close < openClose.Open
}

type Reservation struct {
	Id            string `firestore:"-"`
	RecipientId   string `firestore:"recipientId"`
	Contact       string `firestore:"contact"`
	Name          string `firestore:"name"`
	ScheduledTime int64  `firestore:"scheduledTime"`
	NumPeople     int    `firestore:"numPeople"`
	IsVisible     bool   `firestore:"visible"`
}

func (business Business) TimeClose(day string) int64 {
	return business.Hours[day].CloseTime
}

func (business *Business) IsOpenOnDay(day time.Time) bool {
	dayOfWeek := int(day.Weekday())

	dateString := fmt.Sprinf("%d-%d-%d", day.Year(), day.Month(), day.Day())

	openClose := OpenClose{}
	if val, exists := business.HoursExcpetions[dateString]; exists {
		openClose = val
	} else {
		openClose = business.Hours[dayOfWeek]
	}
	isOpen := openClose.IsOpen

	if !isOpen {
		return false
	}

	currentTimeInt := day.Hour()*100 + day.Minute()

	if openClose.ClosePastMidnight() {
		return currentTimeInt >= openClose.Open || currentTimeInt <= openClose.Close
	} else {
		return openClose.Open <= currentTimeInt <= openClose.Close
	}

}

func (business *Business) IsOpenOnDateTime(dateTime string) bool {
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
	Start          int         `json:"start"`
	End            int         `json:"end"`
	Value          interface{} `json:"value"`
	AdditionalInfo interface{} `json:"additional_info"`
	Text           string      `json:"text"`
	Confidence     float64     `json:"confidence"`
	Entity         string      `json:"entity"`
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

type EnvValues struct {
	PizzaPort string `json:"pizza_port"`
	RasaPort  string `json:"rasa_port"`
}

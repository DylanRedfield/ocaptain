package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"os"
	"strconv"
	"time"
)

type Bot struct {
	Client         *firestore.Client
	Ctx            context.Context
	SwiftClient    SwiftClient
	State          string
	DemoCounter    int
	IsDemo         bool
	BusinessCache  map[string]*Business
	RecipientCache map[string]*Recipient
}

func NewBot(ctx context.Context) (*Bot, error) {
	var bot Bot
	args := os.Args

	bot.BusinessCache = make(map[string]*Business)
	bot.RecipientCache = make(map[string]*Recipient)

	var sa option.ClientOption
	if len(args) > 1 {
		state := args[1]

		if state == "prod" {
			bot.State = PROD_STATE
			sa = option.WithCredentialsFile("prod-firebase-config.json")
		} else if state == "dev_local" {
			bot.State = DEV_STATE_LOCAL
			sa = option.WithCredentialsFile("dev-firebase-config.json")
		} else if state == "dev" {
			bot.State = DEV_STATE_ONLINE
			sa = option.WithCredentialsFile("dev-firebase-config.json")
		}

		if len(args) > 2 {
			demoFlag := args[2]

			if demoFlag == "demo" {
				bot.IsDemo = true
			}
		}

	} else {
		bot.State = DEV_STATE_LOCAL
		sa = option.WithCredentialsFile("dev-firebase-config.json")
	}

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

	swiftClient := SwiftClient{
		AccountKey: "8hjeuf40gqyFFkY1wnL7ikTba1zg3fEk"}

	bot.Client = client
	bot.Ctx = ctx
	bot.SwiftClient = swiftClient
	return &bot, nil
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
	Id              string   `firestore:"-"`
	Name            string   `firestore:"name,omitempty"`
	Address         string   `firestore:"name,omitempty"`
	Contact         string   `firestore:"contact"`
	Platform        Platform `firestore:"platform,omitempty"`
	RecentMessage   *Message `firestore:"recentMessage,omitempty"`
	RecentOrderId   string   `firestore:"recentOrderId,omitempty"`
	TimeLastQueried time.Time
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
	Id                               string               `firestore:"-"`
	Approved                         bool                 `firestore:"approved"`
	Password                         string               `firestore:"password"`
	Whatsapp                         string               `firestore:"whatsapp"`
	PhoneNumber                      string               `firestore:"phoneNumber"`
	Hours                            map[string]OpenClose `firestore:"hours"`
	HoursExceptions                  map[string]OpenClose `firestore:"hoursExceptions"`
	ReservationPlatform              string               `firestore:"reservationPlatform"`
	ReservationPlatformId            string               `firestore:"reservationPlatformId"`
	Employees                        []Employee           `firestore:"employees"`
	SmsPlatform                      string               `firestore:"smsPlatform"`
	SmsNotifyEnabled                 bool                 `firestore:"smsNotifyEnabled"`
	FacebookMessengerId              string               `firestore:"facebookMessengerId"`
	FacebookMessengerPageAccessToken string               `firestore:"facebookMessengerPageAccessToken"`
	TwilioAccountSid                 string               `firestore:"twilioAccountSid"`
	TwilioAuthToken                  string               `firestore:"twilioAuthToken"`
	TwilioClient                     TwilioClient
	FacebookMessengerClient          FacebookMessengerClient
	TimeLastQueried                  time.Time
}

type OpenClose struct {
	IsOpen bool  `firestore:"isOpen"`
	Open   int64 `firestore:"open"`
	Close  int64 `firestore:"close"`
}

type Employee struct {
	IsActive    bool   `firestore:"active"`
	Id          string `firestore:"-"`
	PhoneNumber string `firestore:"phoneNumber"`
}

func (openClose *OpenClose) ClosePastMidnight() bool {
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

func (business *Business) GetOpenCloseOnDay(day time.Time) OpenClose {
	dayOfWeek := int(day.Weekday())
	log.Println(dayOfWeek)

	dateString := fmt.Sprintf("%d-%d-%d", day.Year(), day.Month(), day.Day())

	openClose := OpenClose{}
	if val, exists := business.HoursExceptions[dateString]; exists {
		openClose = val
	} else {
		openClose = business.Hours[strconv.Itoa(dayOfWeek)]
	}

	return openClose

}

func (business *Business) GetNextOpenDayAfter(day time.Time) time.Time {
	// Find the next available open time by add a day at a time in a while loop until the business is open

	// First check if the current day is the next open day by checking if the requested time is less than the open time

	requestedTimeInt := int64(day.Hour()*100 + day.Minute())

	if requestedTimeInt < business.GetOpenCloseOnDay(day).Open {
		return day
	}

	day = day.Add(time.Hour * 24)

	for !business.GetOpenCloseOnDay(day).IsOpen {
		day = day.Add(time.Hour * 24)
	}

	return day
}

func formatIntTimeTwelveHourString(inputTime int64) string {
	period := "am"

	if inputTime >= 1200 {
		period = "pm"
		inputTime = inputTime - 1200
	}

	if inputTime == 0 {
		inputTime = inputTime + 1200
	}

	minutes := inputTime % 100
	hour := inputTime / 100

	return fmt.Sprintf("%d:%02d %s", hour, minutes, period)

}

// Will return error is the business is not open that day
func (business *Business) TimeCloseOnDayString(day time.Time) (string, error) {
	openClose := business.GetOpenCloseOnDay(day)

	if !openClose.IsOpen {
		return "", errors.New("Restaurant closed")
	}

	return formatIntTimeTwelveHourString(openClose.Close), nil
}

func (business *Business) TimeOpenOnDayString(day time.Time) (string, error) {
	openClose := business.GetOpenCloseOnDay(day)

	if !openClose.IsOpen {
		return "", errors.New("Restaurant closed")
	}

	return formatIntTimeTwelveHourString(openClose.Open), nil

}

func (business *Business) IsOpenOnDay(day time.Time) bool {
	openClose := business.GetOpenCloseOnDay(day)
	isOpen := openClose.IsOpen

	if !isOpen {
		log.Println("fuck")
		return false
	}

	currentTimeInt := int64(day.Hour()*100 + day.Minute())

	if openClose.ClosePastMidnight() {
		log.Println("fuck2")
		return currentTimeInt >= openClose.Open || currentTimeInt <= openClose.Close
	} else {
		log.Println(currentTimeInt)
		return openClose.Open <= currentTimeInt && currentTimeInt <= openClose.Close
	}

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
	PizzaPort       string `json:"pizza_port"`
	RasaPort        string `json:"rasa_port"`
	MessengerVerify string `json:"messenger_verify"`
}

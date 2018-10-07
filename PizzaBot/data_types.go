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
	Error       *BotError
}

type BotError struct {
	Type                     BotErrorType
	Message                  string
	ShouldDisplayErrorInChat bool
}

func (err *BotError) Error() string {
	return err.Message
}

type BotErrorType string

const (
	Connection  BotErrorType = "Connection"
	Firebase                 = "Firebase"
	DialogFlow               = "DialogFlow"
	Application              = "Application"
)

// Business
const (
	PhoneNumber string = "phoneNumber"
)

// Recipient
const (
	Contact       string = "contact"
	RecentMessage        = "recentMessage"
	RecentOrderId        = "recentOrderId"
)

// Order
const (
	Type     string = "type"
	Contents        = "contents"
	Name            = "name"
	Address         = "address"
)

// Collections
const (
	Businesses string = "businesses"
	Recipients        = "recipients"
	Messages          = "messages"
	Orders            = "orders"
)

// Actions
const (
	UTTER_GREET                     string = "utter_greet"
	UTTER_GOODBYE                          = "utter_goodbye"
	UTTER_YOUR_WELCOME                     = "utter_your_welcome"
	UTTER_ASK_ADDRESS                      = "utter_ask_address"
	UTTER_ASK_NAME                         = "utter_ask_name"
	UTTER_THANK                            = "utter_thank"
	UTTER_ASK_ORDER_CONTENTS               = "utter_ask_order_contents"
	UTTER_ASK_CONFIRMATION_DELIVERY        = "utter_ask_confirmation_delivery"
	UTTER_ASK_CONFIRMATION_PICK_UP         = "utter_ask_confirmation_pick_up"
	UTTER_ASK_TYPE                         = "utter_ask_type"
	UTTER_AFTER_ORDER                      = "utter_after_order"
  UTTER_ASK_IS_ALL = "utter_ask_is_all"
	ACTION_LISTEN                          = "action_listen"
	ACTION_START_ORDER                     = "action_start_order"
	ACTION_START_ORDER_WITH_INPUTS          = "action_start_order_with_inputs"
	ACTION_SET_TYPE                        = "action_set_type"
	ACTION_SET_ADDRESS                     = "action_set_address"
	ACTION_SET_CONTENT                     = "action_set_content"
	ACTION_SET_NAME                        = "action_set_name"
	ACTION_CHECK_IS_OPEN                   = "action_check_is_open"
	ACTION_CHECK_IS_OPEN_ON_DAY            = "action_check_is_open_on_day"
	ACTION_CHECK_TIME_CLOSE                = "action_check_time_close"
	ACTION_CHECK_TIME_CLOSE_ON_DAY         = "action_check_time_close_on_DAY"
  ACTION_UPDATE_ORDER = "action_update_order"
  ACTION_RESET_SLOTS = "action_reset_slots"
)

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
	Id        string `json:"id"`
	Recipient *Recipient `json:"recipient"`
	Message   *Message `json:"message"`
	Business  *Business `json:"business"`
}

type OutsideResponse struct {
}

func NewBotError(message string, errorType BotErrorType, shouldDisplayErrorInChat bool) BotError {
	return BotError{Type: errorType, Message: message, ShouldDisplayErrorInChat: shouldDisplayErrorInChat}
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
	Id          string `firestore:"-"`
	Approved    bool   `firestore:"approved"`
	Password    string `firestore:"password"`
	PhoneNumber string `firestore:"phoneNumber"`
}

func (business Business) TimeClose() string {
  // TODO implelemt and stop returning string
  return "9:00pm"
}

type Tracker struct {
	Slots    map[string]string `json:"slots"`
	SenderId string            `json:"sender_id"`
  LatestMessage LatestMessage `json:"lastest_message"`
}

type LatestMessage struct {
  Text string `json:"text"`
  Intent string `json:"intent"`
  Entities []Entity `json:"entities"`
}

type Entity struct {
  Start int32 `json:"start"`
  End int32 `json:"end"`
  Value string `json:"value"`
  Text string `json:"text"`
  Confidence float64 `json:"confidence"`
  Entity string `json:"entity"`
}

type RasaRequest struct {
	NextAction string  `json:"next_action"`
  SenderId string `json:"sender_id"`
	Tracker    Tracker `json:"tracker"`
}

type RasaResponse struct {
  Events []Event `json:"events"`
  Responses []Response `json:"responses"`
}
func NewRasaResponse() *RasaResponse {
  return &RasaResponse{
    Events: []Event{},
    Responses: []Response{},
  }
}

type Response struct {
  Text string `json:"text"`
}

type Event struct {
  Event string `json:"event"`
  Name string `json:"name"`
  Value []map[string]string
}

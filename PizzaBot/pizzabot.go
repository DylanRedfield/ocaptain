package main

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Bot struct {
	Client    *firestore.Client
	Ctx       context.Context
	SmsClient TwilioClient
}


func NewBot(ctx context.Context) (Bot, *BotError) {
	sa := option.WithCredentialsFile("firebase-config.json")

	app, err := firebase.NewApp(ctx, nil, sa)

	if err != nil {
		log.Println(err)
	}

	client, err := app.Firestore(ctx)

	if err != nil {
		log.Println(err)
	}

	twilioClient := TwilioClient{
		AccountSid: "AC9dfbda388f3ee10353bbc001694f5c27",
		AuthToken:  "e3429e06cc27740f1c859d2bfc9964ae"}

	return Bot{Client: client, Ctx: ctx, SmsClient: twilioClient}, nil
}

var bot Bot

func init() {
	ctx = context.Background()

  var err *BotError
  bot, err = NewBot(ctx)

  if err != nil {
    log.Println("Error initiating bot")
  }
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

type Tracker struct {
	Slots    map[string]string `json:"slots"`
	SenderId string            `json:"sender_id"`
	// LatestMessage LatestMessage
}

type RasaResponse struct {
	NextAction string  `json:"next_action"`
  SenderId string `json:"sender_id"`
	Tracker    Tracker `json:"tracker"`
}

// Recieves Botrequest, saves message to firebase, sends to recipient, and returns reponse
func HandleBusinessInput(ctx context.Context, reqObj BusinessRequest) BusinessResponse {
	messageRef := bot.Client.Collection(Businesses).Doc(reqObj.Business.Id).Collection(Messages).NewDoc()

	timeInMil := time.Now().UnixNano() / 1000000

	message := Message{
		Content:          reqObj.Message,
		IsBusinessSender: true,
		TimeSent:         timeInMil,
		HasBusinessRead:  true,
		DidBotCreate:     false,
		RecipientId:      reqObj.Recipient.Id}

	_, err := messageRef.Set(ctx, message)

	personRef := bot.Client.Collection(Businesses).Doc(reqObj.BusinessId).Collection(Recipients).Doc(reqObj.Recipient.Id)
	personRef.Update(ctx, []firestore.Update{
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

/*
*actions:
  - utter_greet
  - utter_goodbye
  - utter_your_welcome
  - utter_ask_address
  - utter_ask_name
  - utter_thank
  - utter_ask_order_contents
  - utter_ask_confirmation_pick_up
  - utter_ask_confirmation_delivery
  - utter_ask_type
  - utter_after_order
  - action_place_order
  - action_check_time_close
  - action_check_is_open
  - action_start_order
  - action_start_order_with_inputs
  - action_check_is_open_on_day
  - action_check_time_close_on_day
  - action_set_content
  - action_set_type
  - action_set_name
  - action_set_address
*/
// Recieves outside input, saves to firebase, processes and potentially responds, and returns response
// Only the Phone number and platform field of recipient are garunteed to be set

/*func HandleOutsideInput(ctx context.Context, reqObj OutsideRequest) OutsideResponse {

	rasaUrl := fmt.Sprintf("http://localhost:5005/conversations/%s/respond", reqObj.Recipient.Id)
	body := []byte(fmt.Sprintf(`{"query":"%s"}`, reqObj.Message.Content))

	req, err := http.NewRequest("POST", rasaUrl, bytes.NewBuffer(body))

	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
	}

	var rasaResp RasaResponse

	if err := json.Unmarshal(body, &rasaResp); err != nil {
		log.Println(err)
	}

	nextAction := rasaResp.NextAction

	log.Println(fmt.Sprintf("Next Action %s", nextAction))
	for nextAction != ACTION_LISTEN {
		switch nextAction {
		case UTTER_GREET:
			bot.ActionUtterGreet(reqObj)
		case UTTER_GOODBYE:
			bot.ActionUtterGoodbye(reqObj)
		case UTTER_YOUR_WELCOME:
			bot.ActionUtterYourWelcome(reqObj)
		case UTTER_ASK_ADDRESS:
			bot.ActionUtterAskAddress(reqObj)
		case UTTER_ASK_NAME:
			bot.ActionUtterAskName(reqObj)
		case UTTER_THANK:
			bot.ActionUtterThank(reqObj)
		case UTTER_ASK_ORDER_CONTENTS:
			bot.ActionUtterAskOrderContents(reqObj)
		case UTTER_ASK_CONFIRMATION_PICK_UP:
			bot.ActionUtterAskConfirmationPickUp(reqObj, rasaResp)
		case UTTER_ASK_CONFIRMATION_DELIVERY:
			bot.ActionUtterAskConfirmationDelivery(reqObj, rasaResp)
		case UTTER_ASK_TYPE:
			bot.ActionUtterAskType(reqObj)
		case UTTER_AFTER_ORDER:
			bot.ActionUtterAfterOrder(reqObj)
    case UTTER_ASK_IS_ALL:
      bot.ActionUtterAskIsAll(reqObj)
		case ACTION_START_ORDER:
			//ActionStartOrder(reqObj)
		case ACTION_START_ORDER_WITH_INPUTS:
			bot.ActionStartOrderWithInputs(reqObj, rasaResp)
    case ACTION_UPDATE_ORDER:
      bot.ActionUpdateOrder(reqObj, rasaResp)
		case ACTION_SET_NAME:
			bot.ActionSetName(reqObj, rasaResp)
		case ACTION_SET_TYPE:
			bot.ActionSetType(reqObj, rasaResp)
		case ACTION_SET_ADDRESS:
			bot.ActionSetAddress(reqObj, rasaResp)
		case ACTION_SET_CONTENT:
			bot.ActionSetContent(reqObj, rasaResp)
    case ACTION_CHECK_IS_OPEN:
      bot.actionUtter(reqObj, "Yes")
    case ACTION_CHECK_IS_OPEN_ON_DAY:
      bot.actionUtter(reqObj, "Yes")
    case ACTION_CHECK_TIME_CLOSE:
      bot.actionUtter(reqObj, "Yes")
    case ACTION_CHECK_TIME_CLOSE_ON_DAY:
      bot.actionUtter(reqObj, "Yes")
    case ACTION_RESET_SLOTS:
      bot.ActionResetSlots(reqObj)
		}

		rasaUrl := fmt.Sprintf("http://localhost:5005/conversations/%s/continue", reqObj.Recipient.Id)
		body := []byte(fmt.Sprintf(`{"executed_action":"%s"}`, rasaResp.NextAction))

		req, err := http.NewRequest("POST", rasaUrl, bytes.NewBuffer(body))

		if err != nil {
			log.Println(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println(err)
		}

		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Println(err)
		}

		if err := json.Unmarshal(body, &rasaResp); err != nil {
			log.Println(err)
		}

		nextAction = rasaResp.NextAction

	}
	return OutsideResponse{}
} */
func HandleAction(req *RasaResponse) {
	action := req.NextAction
  log.Println(action)
		switch action {
		case ACTION_START_ORDER:
			ActionStartOrder(req)
		/*case ACTION_START_ORDER_WITH_INPUTS:
			bot.ActionStartOrderWithInputs(reqObj, rasaResp)
    case ACTION_UPDATE_ORDER:
      bot.ActionUpdateOrder(reqObj, rasaResp)
    case ACTION_CHECK_IS_OPEN:
      bot.actionUtter(reqObj, "Yes")
    case ACTION_CHECK_IS_OPEN_ON_DAY:
      bot.actionUtter(reqObj, "Yes")
    case ACTION_CHECK_TIME_CLOSE:
      bot.actionUtter(reqObj, "Yes")
    case ACTION_CHECK_TIME_CLOSE_ON_DAY:
      bot.actionUtter(reqObj, "Yes")
    case ACTION_RESET_SLOTS:
      bot.ActionResetSlots(reqObj)*/
		}

}
func HandleOutsideInput(ctx context.Context, reqObj OutsideRequest) OutsideResponse {
  // Send a http request that will be handled in the textual_input_channel
  // The body is the OutsideRequest object

	rasaUrl := fmt.Sprintf("http://localhost:5005/webhooks/textual/webhook")
	body, err := json.Marshal(reqObj)

  if err != nil {
    log.Println(err)
  }


	req, err := http.NewRequest("POST", rasaUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Println(err)
	}

	http.DefaultClient.Do(req)

  return OutsideResponse{}
}

func ActionStartOrder(req *RasaResponse) {
	order := Order{
		RecipientId:          req.SenderId,
    RecipientContact: req.Tracker.Slots["recipient_contact"],
		StartTime:            currentTime(),
		LastModificationTime: currentTime(),
		IsVisible:            true,
	}

	slots := req.Tracker.Slots
	if slots["address"] != "" {
		order.Address = slots["address"]
	}

	if slots["name"] != "" {
		order.Name = slots["name"]
	}

	if slots["type"] != "" {
		order.Type = slots["type"]
	}

	if slots["content"] != "" {
		order.Content = slots["content"]
	}

	bot.saveOrder(req, &order)
}

func (bot Bot) ActionStartOrderWithInputs(req OutsideRequest, resp RasaResponse) {
	slots := resp.Tracker.Slots

	order := Order{
		RecipientId:          req.Recipient.Id,
    RecipientContact: req.Recipient.Contact,
		StartTime:            currentTime(),
		LastModificationTime: currentTime(),
		IsVisible:            true,
	}

	if slots["address"] != "" {
		order.Address = slots["address"]
	}

	if slots["name"] != "" {
		order.Name = slots["name"]
	}

	if slots["type"] != "" {
		order.Type = slots["type"]
	}

	if slots["content"] != "" {
		order.Content = slots["content"]
	}

	//bot.saveOrder(req, &order)
}

func (bot Bot) saveOrder(req *RasaResponse, order *Order) {
  businessId := req.Tracker.Slots["business_id"]
  recipientId := req.Tracker.Slots["recipient_id"]

  log.Println(businessId)

	ordersRef := bot.Client.Collection(Businesses).Doc(businessId).Collection(Orders)

	docRef, _, err := ordersRef.Add(bot.Ctx, order)

	if err != nil {
		log.Println(err)
	}

	order.Id = docRef.ID

	recipientRef := bot.Client.Collection(Businesses).Doc(businessId).Collection(Recipients).Doc(recipientId)
	recipientRef.Update(ctx, []firestore.Update{
		{Path: RecentOrderId, Value: order.Id},
	})

}

func (bot Bot) ActionSetAddress(req OutsideRequest, resp RasaResponse) {
	orderRef := bot.Client.Collection(Businesses).Doc(req.Business.Id).Collection(Orders).Doc(req.Recipient.RecentOrderId)
	orderRef.Update(ctx, []firestore.Update{
		{Path: Address, Value: resp.Tracker.Slots["address"]},
	})
}

func (bot Bot) ActionSetType(req OutsideRequest, resp RasaResponse) {
	orderRef := bot.Client.Collection(Businesses).Doc(req.Business.Id).Collection(Orders).Doc(req.Recipient.RecentOrderId)
	orderRef.Update(ctx, []firestore.Update{
		{Path: Type, Value: resp.Tracker.Slots["type"]},
	})
}

func (bot Bot) ActionSetName(req OutsideRequest, resp RasaResponse) {
	orderRef := bot.Client.Collection(Businesses).Doc(req.Business.Id).Collection(Orders).Doc(req.Recipient.RecentOrderId)
	orderRef.Update(ctx, []firestore.Update{
		{Path: Name, Value: resp.Tracker.Slots["name"]},
	})
}

func (bot Bot) ActionSetContent(req OutsideRequest, resp RasaResponse) {
	orderRef := bot.Client.Collection(Businesses).Doc(req.Business.Id).Collection(Orders).Doc(req.Recipient.RecentOrderId)
	orderRef.Update(ctx, []firestore.Update{
		{Path: Contents, Value: resp.Tracker.Slots["content"]},
	})
}

func (bot Bot) ActionUpdateOrder(req OutsideRequest, resp RasaResponse) {
	orderRef := bot.Client.Collection(Businesses).Doc(req.Business.Id).Collection(Orders).Doc(req.Recipient.RecentOrderId)
	orderRef.Update(ctx, []firestore.Update{
		{Path: Contents, Value: resp.Tracker.Slots["content"]},
		{Path: Name, Value: resp.Tracker.Slots["name"]},
		{Path: Address, Value: resp.Tracker.Slots["address"]},
		{Path: Type, Value: resp.Tracker.Slots["type"]},
	})

}

func (bot Bot) ActionResetSlots(req OutsideRequest) {
	rasaUrl := fmt.Sprintf("http://localhost:5005/conversations/%s/tracker/events", req.Recipient.Id)
	body, err := json.Marshal(req)

  if err != nil {
    log.Println(err)
  }

	_, err = http.NewRequest("POST", rasaUrl, bytes.NewBuffer(body))

  if err != nil {
    log.Println(err)
  }

}

func (bot Bot) actionUtter(reqObj OutsideRequest, utterance string) {

	message := Message{
		Content:          utterance,
		IsBusinessSender: true,
		TimeSent:         currentTime(),
		DidBotCreate:     true,
		RecipientId:      reqObj.Recipient.Id,
		HasBusinessRead:  false,
	}

	bot.saveSms(reqObj.Recipient, *reqObj.Business, &message)

	smsReq := SMSRequest{
		To:   reqObj.Recipient.Contact,
		From: reqObj.Business.PhoneNumber,
		Body: utterance,
	}

	bot.SmsClient.SendSMS(smsReq)
}

func currentTime() int64 {
	return time.Now().UnixNano() / 1000000
}

func randomItem(choices []string) string {
	randomIndex := rand.Int31n(int32(len(choices)))

	return choices[randomIndex]
}

func (bot Bot) saveSms(recipient *Recipient, business Business, message *Message) {
	messagesRef := bot.Client.Collection(Businesses).Doc(business.Id).Collection(Messages)

	docRef, _, _ := messagesRef.Add(bot.Ctx, message)

	message.Id = docRef.ID

	log.Println(fmt.Sprintf("Recipient: %s business: %s message %s", recipient.Id, business.Id, message.Id))

	recipientRef := bot.Client.Collection(Businesses).Doc(business.Id).Collection(Recipients).Doc(recipient.Id)
	recipientRef.Update(ctx, []firestore.Update{
		{Path: RecentMessage, Value: message},
	})

	recipient.RecentMessage = message

}

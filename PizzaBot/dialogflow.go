package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

const clientToken = "e5490e090a2442a085bb73466a84237d"

var baseUrl string = "https://api.dialogflow.com/v1/"
var client = &http.Client{Timeout: time.Second * 10}

type CreateEntriesData struct {
	Name    string  `json:"name"`
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Value    string   `json:"value"`
	Synonyms []string `json:"synonyms"`
}

type ContextParameter struct {
	Name     string
	Value    string
	Original string
}

type Input struct {
	Query     string
	SessionId string
}

func init() {
	log.SetFlags(log.Lshortfile)
}

// Need to implement Marshler because the key is dynamic
/*func (p Context) MarshalJSON() ([]byte, error) {
  buffer := bytes.NewBufferString("{")
  buffer.WriteString(fmt.Sprintf("\"%v\":[\"%v\"],", p.Name, strings.Join(p.Values, "\",\"")))
  buffer.WriteString(fmt.Sprintf("\"%v.original\":\"%v\",", p.Name, p.Original))
  buffer.WriteString(fmt.Sprintf("\"intent_action\":\"%v\"}", p.Action))
  return buffer.Bytes(), nil
} */

func (p *ContextParameter) UnmarshalJSON(rawBytes []byte) error {
	var asMap map[string]string

	if err := json.Unmarshal(rawBytes, &asMap); err != nil {
		return err
	}

	for key, value := range asMap {
		reg := regexp.MustCompile("original")
		containsOriginal := reg.MatchString(key)

		if !containsOriginal {
			p.Name = key
			p.Value = value
		} else {
			p.Original = value
		}
	}

	return nil
}

type Context struct {
	Lifespan   int              `json:"lifespan,omitempty"`
	Name       string           `json:"name,omitempty"`
	Parameters ContextParameter `json:"parameters,omitempty"`
}

type FlowRequest struct {
	Contexts      []Context `json:"contexts,omitempty"`
	Query         string    `json:"query,omitempty"`
	ResetContexts bool      `json:"resetContexts,omitempty"`
	SessionId     string    `json:"sessionId,omitempty"`
	Lang          string    `json:"lang"`
}

type FlowMessage struct {
	Speech string
	Type   int
}
type FlowFulfillment struct {
	Messages []FlowMessage
	Speech   string
}
type FlowResultMetadata struct {
	IntentId                  string
	IntentName                string
	WebhookForSlotFillingUsed string
	WebhookResponseTime       int
	WebhookUsed               string
}
type FlowResult struct {
	Action           string
	ActionIncomplete bool
	Contexts         []Context
	Fulfillment      FlowFulfillment
	Metadata         FlowResultMetadata
	Parameters       map[string]string
	ResolvedQuery    string
	Score            float32
	Source           string
}

type FlowStatus struct {
	Code                    int
	Eis_runningrrrorDetails string
	ErrorType               string
}
type FlowResponse struct {
	Id        string
	Lang      string
	Result    FlowResult
	SessionId string
	Status    FlowStatus
	Timestamp string
}

type BotRequest struct {
	Query     string
	SessionId string
	Lat, Lon  float64 `json:",string"`
}

type BotResponse struct {
	Message string
	Error   *BotError `json:",omitempty"`
}

func Query(input Input) (*FlowResponse, *BotError) {
	data, err := json.Marshal(FlowRequest{Lang: "en", Query: input.Query, SessionId: input.SessionId})

	if err != nil {
		log.Println(err)
		return nil, &BotError{Message: err.Error(), Type: Application}
	}

	bufferedJson := bytes.NewBuffer(data)

	url := "https://api.dialogflow.com/v1/query?v=20170712"

	request, err := http.NewRequest("POST", url, bufferedJson)

	request.Header.Add("Authorization", "Bearer e5490e090a2442a085bb73466a84237d")
	request.Header.Add("Content-Type", "application/json")

	rawResponse, err := client.Do(request)
	defer rawResponse.Body.Close()

	if err != nil {
		log.Println(err)
		return nil, &BotError{Message: err.Error(), Type: Application}
	}

	body, err := ioutil.ReadAll(rawResponse.Body)

	if err != nil {
		log.Println(err)
		return nil, &BotError{Message: err.Error(), Type: Application}
	}

	var response FlowResponse

	if err := json.Unmarshal(body, &response); err != nil {
		log.Println(err)
		return nil, &BotError{Message: err.Error(), Type: Application}
	}

	log.Println(response)
	return &response, nil
}

func Misunderstand() string {
	phrases := []string{
		"Sorry, I didn't quite get that, can you try again?",
		"I can't understand your phrasing, can you reword it for me?",
		"I can't understand that. Try again."}

	return phrases[rand.Int31n(int32(len(phrases)))%int32(len(phrases))]
}

/*    switch err.Type {
      message = RandomNetworkError()
      message = RandomNextbusInputError()
      message = RandomDialogFlowInputError()
      message = RandomApplicationError()
    } */
func RandomNetworkError() string {
	phrases := []string{
		"I'm sorry I'm having trouble connecting to the internet right now",
		"Im sorry but I'm have some network issues right now",
		"I'm sorry but my internet isn't working"}

	return phrases[rand.Int31n(int32(len(phrases)))%int32(len(phrases))]
}

func RandomNextbusError() string {
	phrases := []string{
		"Sorry, but NextBus isn't working right now, sorry",
		"Sorry, NextBus is down",
		"Sorry, NextBus is down, so I'm pretty useless right now"}

	return phrases[rand.Int31n(int32(len(phrases)))%int32(len(phrases))]
}

func RandomDialogflowError() string {
	phrases := []string{
		"Placeholder",
		"Placeholder",
		"Placeholder"}

	return phrases[rand.Int31n(int32(len(phrases)))%int32(len(phrases))]
}

func RandomApplicationError() string {
	phrases := []string{
		"I'm having trouble processing that input. Reword your response and try again",
		"That input confuses me, maybe reword it and try again?",
		"I can't understand that input, please reword it"}

	return phrases[rand.Int31n(int32(len(phrases)))%int32(len(phrases))]
}

func RandomNoLocationError() string {
	phrases := []string{
		"I couldn't get your location, so please make sure to specifify a stop",
		"Without your location I can't figure that out unless you specifiy a stop",
		"I can't answer that without your location or a specified stop"}

	return phrases[rand.Int31n(int32(len(phrases)))%int32(len(phrases))]
}
func RandomNoIntentMatched() string {
	phrases := []string{
		"Unless I'm misunderstanding your question, I don't know how to do that",
		"I might not understand you, but I don't think I know how to do that"}

	return phrases[rand.Int31n(int32(len(phrases)))%int32(len(phrases))]
}

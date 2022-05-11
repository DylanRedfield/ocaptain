package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type FacebookMessengerClient struct {
	PageAccessToken string
}
type MessengerWebhook struct {
	Entry  []Entry `json:entry`
	Object string  `json:"object"`
}
type Entry struct {
	Messaging []FacebookMessengerReceiveMessage `json:"messaging"`
}
type FacebookMessengerReceiveMessage struct {
	Sender    FacebookSender    `json:"sender"`
	Recipient FacebookRecipient `json:"recipient"`
	Message   FacebookMessage   `json:"message"`
}

type FacebookMessengerSendMessage struct {
	Recipient FacebookRecipient `json:"recipient"`
	Message   FacebookMessage   `json:"message"`
}

type FacebookSender struct {
	Id string `json:"id"`
}
type FacebookRecipient struct {
	Id string `json:"id"`
}

type FacebookMessage struct {
	Text string `json:"text"`
}

func (client *FacebookMessengerClient) Send(req *MessageRequest) {
	var reqObj FacebookMessengerSendMessage

	reqObj.Recipient = FacebookRecipient{Id: req.To}
	reqObj.Message = FacebookMessage{Text: req.Body}

	body, err := json.Marshal(reqObj)

	if err != nil {
		log.Println(err)
	}

	// TODO if using facebook messenger we need each business's access token
	resp, err := http.Post("https://graph.facebook.com/v13.0/me/messages?access_token="+client.PageAccessToken,
		"application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

}

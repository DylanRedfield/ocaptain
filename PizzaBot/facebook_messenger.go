package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type MessengerWebhook struct {
	entry []Entry
}
type Entry struct {
	messaging []FacebookMessengerReceiveMessage
}
type FacebookMessengerReceiveMessage struct {
	sender    FacebookSender
	recipient FacebookRecipient
	message   FacebookMessage
}

type FacebookMessengerSendMessage struct {
	recipient FacebookRecipient
	message   FacebookMessage
}

type FacebookSender struct {
	id string
}
type FacebookRecipient struct {
	id string
}

type FacebookMessage struct {
	text string
}

func Send(req *MessageRequest) {
	var reqObj FacebookMessengerSendMessage

	reqObj.recipient = FacebookRecipient{id: req.To}
	reqObj.message = FacebookMessage{text: req.Body}

	bussiness, err := businessFromFacebookId(req.From)

	if err != nil {
		log.Println(err)
	}

	body, err := json.Marshal(reqObj)

	// TODO if using facebook messenger we need each business's access token
	resp, err := http.Post("https://graph.facebook.com/v13.0/me/messages?access_token="+bussiness.FacebookMessengerPageAccessToken,
		"application/json", bytes.NewBuffer(body))

	defer resp.Body.Close()

}

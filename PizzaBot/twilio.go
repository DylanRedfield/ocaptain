package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	//"io/ioutil"
	"strings"
)

type TwilioClient struct {
	AccountSid string
	AuthToken  string
}

func init() {
	log.SetFlags(log.Lshortfile)
}

func (client *TwilioClient) Send(reqObj *MessageRequest) {
	queryUrl := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", client.AccountSid)

	body := url.Values{}
	body.Set("To", reqObj.To)
	body.Set("From", reqObj.From)
	body.Set("Body", reqObj.Body)

	req, err := http.NewRequest("POST", queryUrl, strings.NewReader(body.Encode()))

	if err != nil {
		// TODO handle error
		log.Println(err)
	}

	req.SetBasicAuth(client.AccountSid, client.AuthToken)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		// TODO handlje error
		log.Println(err)
	}
	defer resp.Body.Close()

	//fuck, err := ioutil.ReadAll(resp.Body)
	// TODO handle response errors
}

func (client *TwilioClient) SendBulk(reqObj *BulkMessageRequest) {

	single := &MessageRequest{From: GetEnvValues().TwilioGeneralNumber, Body: reqObj.Body}
	for _, number := range reqObj.To {
		single.To = number
		client.Send(single)
	}

}

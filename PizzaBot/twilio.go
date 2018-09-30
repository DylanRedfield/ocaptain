package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	//"io/ioutil"
	"strings"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

type TwilioClient struct {
	AccountSid string
	AuthToken  string
}

type SMSRequest struct {
	To   string
	From string
	Body string
}

type TwilioRequest struct {
	Body string
	To   string
	From string
}

func init() {
	log.SetFlags(log.Lshortfile)
}

func (client TwilioClient) SendSMS(reqObj SMSRequest) {

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

	defer resp.Body.Close()

	if err != nil {
		// TODO handlje error
		log.Println(err)
	}

	//fuck, err := ioutil.ReadAll(resp.Body)
	// TODO handle response errors
}

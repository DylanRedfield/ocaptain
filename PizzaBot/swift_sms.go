package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type SwiftClient struct {
	AccountKey string
}

func (client *SwiftClient) Send(req *MessageRequest) {
	baseUrl := fmt.Sprint("http://smsgateway.ca/services/message.svc/%s/%s/ViaDedicated", client.AccountKey, req.To)
	body := url.Values{}
	body.Set("MessageBody", req.Body)
	body.Set("Reference", "")
	body.Set("SenderNumber", req.From)

	httpReq, err := http.NewRequest("POST", baseUrl, strings.NewReader(body.Encode()))
	if err != nil {
		log.Println(err)
	}

	httpReq.Header.Add("Accept", "application/json")
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		// TODO handlje error
		log.Println(err)
	}
	defer resp.Body.Close()
}
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
  "bytes"
	"strings"
  "encoding/json"
)

type SwiftClient struct {
	AccountKey string
}

func (client *SwiftClient) Send(req *MessageRequest) {
  log.Println(req.From)
  log.Println(client.AccountKey)

  baseUrl := fmt.Sprintf("http://smsgateway.ca/services/message.svc/%s/%s/ViaDedicated", client.AccountKey, req.To[1:])

  requestBody, err := json.Marshal(map[string]string{
    "MessageBody": req.Body,
    "SenderNumber": req.From[1:],
    "Reference": "1",
  })

  if err != nil {
    log.Println(err)
  }

	resp, err := http.Post(baseUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		// TODO handlje error
		log.Println(err)
	}
  log.Println(resp.Status)

  buf := new(bytes.Buffer)
  buf.ReadFrom(resp.Body)
  log.Println(buf.String())
	defer resp.Body.Close()
}

func (client *SwiftClient) SendBulk(req *BulkMessageRequest) {
	baseUrl := fmt.Sprintf("http://smsgateway.ca/services/message.svc/%s/Bulk", client.AccountKey)


	numbersToString := fmt.Sprintf("[%s]", strings.Join(req.To, ","))
	body.Set("CellNumbers", numbersToString)

  requestBody, err := json.Marshal(map[string]string{
    "MessageBody": req.Body,
    "Reference": "1",
    "CellNumbers": numbersToString,
  })

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

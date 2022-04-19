package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type SwiftClient struct {
	AccountKey string
}

func (client *SwiftClient) Send(req *MessageRequest) {
	log.Println(req.From)
	log.Println(client.AccountKey)

	baseUrl := fmt.Sprintf("http://smsgateway.ca/services/message.svc/%s/%s/ViaDedicated", client.AccountKey, req.To[1:])

	requestBody, err := json.Marshal(map[string]string{
		"MessageBody":  req.Body,
		"SenderNumber": req.From[1:],
		"Reference":    "1",
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

	if len(req.To) < 1 {
		return
	}

	requestBody, err := json.Marshal(map[string]string{
		"MessageBody": req.Body,
		"Reference":   "1",
		"CellNumbers": numbersToString(req.To),
	})

  log.Println(numbersToString(req.To))

	resp, err := http.Post(baseUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println(err)
	}
	log.Println(resp.Status)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	log.Println(buf.String())

	defer resp.Body.Close()
}

func numbersToString(numbers []string) string {

	noPlus := []string{}

	for _, number := range numbers {
		noPlus = append(noPlus, number[1:])
	}

	return fmt.Sprintf("%s", fmt.Sprintf("\"%s\"", strings.Join(noPlus, ",")))

}

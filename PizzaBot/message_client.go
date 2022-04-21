package main

import (
	"net/http"
	"time"
)

type MessageRequest struct {
	To       string
	From     string
	Body     string
	Platform string
}

type BulkMessageRequest struct {
	To   []string
	Body string
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

type MessageClient interface {
	Send(req *MessageRequest)
}

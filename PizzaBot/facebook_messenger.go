package main

type FacebookMessengerTextMessage struct {
	sender    FacebookSender
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

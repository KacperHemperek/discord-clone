package models

type Message struct {
	Text string `json:"text"`
	Base
}

type MessageWithUser struct {
	User *User `json:"user"`
	Message
}

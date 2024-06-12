package models

import (
	"github.com/kacperhemperek/discord-go/types"
	"time"
)

type Chat struct {
	Name string         `json:"name"`
	Type types.ChatType `json:"type"`
	Base
}

type UserToChat struct {
	ChatID    string    `json:"chatId"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ChatWithMembers struct {
	Members []*User `json:"members"`
	Chat
}

type ChatWithMessages struct {
	Messages []*MessageWithUser `json:"messages"`
	Chat
}

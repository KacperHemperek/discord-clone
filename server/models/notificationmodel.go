package models

import (
	"github.com/kacperhemperek/discord-go/types"
)

type BaseNotification = struct {
	Base
	Type   types.NotificationType `json:"type"`
	Seen   bool                   `json:"seen"`
	UserID int                    `json:"userId"`
}

type NotificationDTO = struct {
	Base
	BaseNotification
	Data []byte `json:"-"`
}

type FriendRequestNotificationData struct {
	TestValue string `json:"testValue" validate:"required"`
}

type FriendRequestNotification struct {
	BaseNotification
	Data *FriendRequestNotificationData `json:"data"`
}

type NewMessageNotificationData struct {
	ChatID string `json:"chatId"`
}

type NewMessageNotification struct {
	BaseNotification
	Data NewMessageNotificationData `json:"data"`
}

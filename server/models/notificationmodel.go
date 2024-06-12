package models

import (
	"encoding/json"
	"errors"
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
	ChatID int `json:"chatId"`
}

func (n *NewMessageNotificationData) Scan(value any) error {
	switch val := value.(type) {
	case []byte:
		return json.Unmarshal(val, n)
	case string:
		return json.Unmarshal([]byte(val), n)
	default:
		return errors.New("invalid new messsage notification data")
	}
}

type NewMessageNotification struct {
	BaseNotification
	Data NewMessageNotificationData `json:"data"`
}

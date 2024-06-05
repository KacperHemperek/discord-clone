package models

import (
	"github.com/kacperhemperek/discord-go/types"
)

type NotificationDTO = struct {
	Base
	Type   types.NotificationType `json:"type"`
	Seen   bool                   `json:"seen"`
	UserID int                    `json:"userId"`
	Data   []byte                 `json:"-"`
}

type BaseFriendRequestNotification = struct {
	Base
	Type   types.NotificationType `json:"type"`
	Seen   bool                   `json:"seen"`
	UserID int                    `json:"userId"`
}

type FriendRequestNotificationData = struct {
}

type FriendRequestNotification = struct {
	BaseFriendRequestNotification
	Data FriendRequestNotificationData `json:"data"`
}

type NewMessageNotificationData = struct {
	ChatID string `json:"chatId"`
}

type NewMessageNotification = struct {
	BaseFriendRequestNotification
	Data FriendRequestNotificationData `json:"data"`
}

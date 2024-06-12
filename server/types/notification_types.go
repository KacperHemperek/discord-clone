package types

import (
	"encoding/json"
	"errors"
)

type NotificationType int64

var (
	InvalidNotificationTypeErr = errors.New("invalid notification type")
)

const (
	FriendRequestNotification NotificationType = iota
	NewMessageNotification
)

func (n *NotificationType) String() string {
	switch *n {
	case FriendRequestNotification:
		return "friend_request"
	case NewMessageNotification:
		return "new_message"
	default:
		return "unsupported_notification_type"
	}
}

func (n *NotificationType) UnmarshalJSON(data []byte) error {
	dataStr := string(data)

	switch dataStr {
	case "friend_request":
		*n = FriendRequestNotification
		return nil
	case "new_message":
		*n = NewMessageNotification
		return nil
	default:
		return InvalidNotificationTypeErr
	}
}

func (n *NotificationType) MarshalJSON() ([]byte, error) {
	switch *n {
	case FriendRequestNotification:
		return json.Marshal("friend_request")
	case NewMessageNotification:
		return json.Marshal("new_message")
	default:
		return []byte(""), errors.New("invalid notification type")
	}
}

func (n *NotificationType) Scan(value any) error {
	switch value.(type) {
	case string:
		{
			if value == "friend_request" {
				*n = FriendRequestNotification
				return nil
			}

			if value == "new_message" {
				*n = NewMessageNotification
				return nil
			}

			return InvalidNotificationTypeErr
		}
	default:
		return InvalidNotificationTypeErr
	}
}

func IsNotificationType(value string) bool {
	newMessage := NewMessageNotification
	friendRequest := FriendRequestNotification
	return value == newMessage.String() || value == friendRequest.String()
}

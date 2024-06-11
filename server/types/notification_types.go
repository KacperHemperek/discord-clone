package types

import (
	"encoding/json"
	"errors"
)

type NotificationType int64

const (
	FriendRequestNotification NotificationType = iota
	NewMessageNotification
)

func (n NotificationType) String() string {
	switch n {
	case FriendRequestNotification:
		return "friend_request"
	case NewMessageNotification:
		return "new_message"
	default:
		return "unsupported_notification_type"
	}
}

func (n NotificationType) UnmarshalJSON(data []byte) error {
	dataStr := string(data)

	switch dataStr {
	case "friend_request":
		n = FriendRequestNotification
		break
	case "new_message":
		n = NewMessageNotification
		break
	default:
		return errors.New("invalid notification type")
	}
	return nil
}

func (n NotificationType) MarshalJSON() ([]byte, error) {
	switch n {
	case FriendRequestNotification:
		return json.Marshal("friend_request")
	case NewMessageNotification:
		return json.Marshal("new_message")
	default:
		return []byte(""), errors.New("invalid notification type")
	}
}

func (n NotificationType) Scan(value any) error {
	switch value.(type) {
	case string:
		{
			if value == "friend_request" {
				value = FriendRequestNotification
			}

			if value == "new_message" {
				value = NewMessageNotification
			}

			return nil
		}
	default:
		return errors.New("invalid notification type")
	}
}

func IsNotificationType(value string) bool {
	return value == NewMessageNotification.String() || value == FriendRequestNotification.String()
}

package types

import (
	"encoding/json"
	"errors"
)

type ChatType int64

var (
	InvalidChatTypeErr = errors.New("invalid chat type")
)

const (
	PrivateChat ChatType = iota
	GroupChat
)

func (n *ChatType) String() string {
	switch *n {
	case PrivateChat:
		return "private"
	case GroupChat:
		return "group"
	default:
		return ""
	}
}

func (n *ChatType) UnmarshalJSON(data []byte) error {
	dataStr := string(data)

	switch dataStr {
	case "private":
		*n = PrivateChat
		break
	case "group":
		*n = GroupChat
		break
	default:
		return InvalidChatTypeErr
	}
	return nil
}

func (n *ChatType) MarshalJSON() ([]byte, error) {
	switch *n {
	case PrivateChat:
		return json.Marshal("private")
	case GroupChat:
		return json.Marshal("group")
	default:
		return []byte(""), InvalidChatTypeErr
	}
}

func (n *ChatType) Scan(value any) error {
	switch value.(type) {
	case string:
		{
			if value == "private" {
				*n = PrivateChat
				return nil
			}

			if value == "group" {
				*n = GroupChat
				return nil
			}
			return InvalidChatTypeErr
		}
	default:
		return InvalidChatTypeErr
	}
}

func (n *ChatType) Is(comp ChatType) bool {
	return n.String() == comp.String()
}

package ws

import (
	"github.com/gorilla/websocket"
	"github.com/kacperhemperek/discord-go/models"
)

type ChatServiceInterface interface {
	AddChatConn(chatID, userID int, conn *websocket.Conn)
	BroadcastMessage(chatID int, message *models.MessageWithUser)
	CloseConn(chatID, userID int)
}

type ChatService struct {
	chats map[int]map[int]*websocket.Conn
}

func (s *ChatService) AddChatConn(chatID, userID int, conn *websocket.Conn) {
	_, ok := s.chats[chatID]
	if !ok {
		newConns := map[int]*websocket.Conn{
			userID: conn,
		}
		s.chats[chatID] = newConns
	} else {
		s.chats[chatID][userID] = conn
	}
}

func (s *ChatService) BroadcastMessage(chatID int, message *models.MessageWithUser) {
	nm := &newMessage{
		Type:    NewMessage,
		Message: message,
	}
	conns, ok := s.chats[chatID]
	if !ok {
		return
	}
	for _, conn := range conns {
		err := conn.WriteJSON(nm)
		if err != nil {
			return
		}
	}
}

func (s *ChatService) CloseConn(chatID, userID int) {
	chatConns, chatFound := s.chats[chatID]
	if chatFound {
		conn, connFound := chatConns[userID]
		if connFound {
			err := conn.Close()
			delete(chatConns, userID)
			if err != nil {
				return
			}
		}
	}
}

func NewChatService() *ChatService {
	return &ChatService{
		chats: make(map[int]map[int]*websocket.Conn),
	}
}

type newMessage struct {
	Type    string                  `json:"type"`
	Message *models.MessageWithUser `json:"message"`
}

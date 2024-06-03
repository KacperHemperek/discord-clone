package ws

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/kacperhemperek/discord-go/models"
	"sync"
)

var (
	ChatNotFound = errors.New("chat connection object not found")
)

type ChatServiceInterface interface {
	AddChatConn(chatID, userID int, conn *websocket.Conn)
	BroadcastNewMessage(chatID int, message *models.MessageWithUser) error
	BroadcastNewChatName(chatID int, name string) error
	CloseConn(chatID, userID int) error
}

type ChatService struct {
	chats     map[int]map[int]*websocket.Conn
	chatsLock sync.RWMutex
}

func (s *ChatService) AddChatConn(chatID, userID int, conn *websocket.Conn) {
	s.chatsLock.Lock()
	_, ok := s.chats[chatID]
	defer s.chatsLock.Unlock()
	if !ok {
		newConns := map[int]*websocket.Conn{
			userID: conn,
		}
		s.chats[chatID] = newConns
	} else {
		s.chats[chatID][userID] = conn
	}
}

func (s *ChatService) BroadcastNewMessage(chatID int, message *models.MessageWithUser) error {
	nm := newNewMessage(message)
	return s.broadcastMessage(chatID, nm)
}

func (s *ChatService) BroadcastNewChatName(chatID int, name string) error {
	changeNameMessage := newChatNameChanged(name)
	return s.broadcastMessage(chatID, changeNameMessage)
}

func (s *ChatService) CloseConn(chatID, userID int) error {
	s.chatsLock.Lock()
	defer s.chatsLock.Unlock()
	chatConns, chatFound := s.chats[chatID]
	if chatFound {
		conn, connFound := chatConns[userID]
		if connFound {
			err := conn.Close()
			delete(chatConns, userID)
			if err != nil {
				return err
			}
		}
	}
	return ChatNotFound
}

func (s *ChatService) broadcastMessage(chatID int, message any) error {
	s.chatsLock.Lock()
	defer s.chatsLock.Unlock()
	chatConns, chatFound := s.chats[chatID]
	if !chatFound {
		return ChatNotFound
	}
	return broadcast(message, chatConns)
}

func NewChatService() *ChatService {
	return &ChatService{
		chats:     make(map[int]map[int]*websocket.Conn),
		chatsLock: sync.RWMutex{},
	}
}

func newNewMessage(m *models.MessageWithUser) *newMessage {
	return &newMessage{
		Type:    NewMessage,
		Message: m,
	}
}

func newChatNameChanged(name string) *chatNameChanged {
	return &chatNameChanged{
		Type:    ChatNameUpdated,
		NewName: name,
	}
}

type chatNameChanged struct {
	Type    string `json:"type"`
	NewName string `json:"newName"`
}

type newMessage struct {
	Type    string                  `json:"type"`
	Message *models.MessageWithUser `json:"message"`
}

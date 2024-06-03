package ws

import (
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kacperhemperek/discord-go/models"
	"sync"
)

var (
	ChatNotFound = errors.New("chat connection object not found")
)

type ChatServiceInterface interface {
	AddChatConn(chatID int, conn *websocket.Conn) string
	BroadcastNewMessage(chatID int, message *models.MessageWithUser) error
	BroadcastNewChatName(chatID int, name string) error
	CloseConn(chatID int, connID string) error
}

type ChatService struct {
	chats     map[int]map[string]*websocket.Conn
	chatsLock sync.RWMutex
}

func (s *ChatService) AddChatConn(chatID int, conn *websocket.Conn) string {
	s.chatsLock.Lock()
	_, chatFound := s.chats[chatID]
	defer s.chatsLock.Unlock()
	connID := uuid.New().String()
	if !chatFound {
		newConns := map[string]*websocket.Conn{
			connID: conn,
		}
		s.chats[chatID] = newConns
	} else {
		s.chats[chatID][connID] = conn
	}
	return connID
}

func (s *ChatService) BroadcastNewMessage(chatID int, message *models.MessageWithUser) error {
	nm := newNewMessage(message)
	return s.broadcastMessage(chatID, nm)
}

func (s *ChatService) BroadcastNewChatName(chatID int, name string) error {
	changeNameMessage := newChatNameChanged(name)
	return s.broadcastMessage(chatID, changeNameMessage)
}

func (s *ChatService) CloseConn(chatID int, connID string) error {
	s.chatsLock.Lock()
	defer s.chatsLock.Unlock()
	chatConns, chatFound := s.chats[chatID]
	if chatFound {
		conn, connFound := chatConns[connID]
		if connFound {
			err := conn.Close()
			if err != nil {
				return err
			}
			delete(chatConns, connID)
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
		chats:     make(map[int]map[string]*websocket.Conn),
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

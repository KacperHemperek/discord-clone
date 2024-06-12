package ws

import (
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kacperhemperek/discord-go/models"
	"sync"
)

var (
	ChatNotFoundErr = errors.New("chat connection object not found")
)

type ChatServiceInterface interface {
	AddChatConn(chatID, userID int, conn *websocket.Conn) string
	BroadcastNewMessage(chatID int, message *models.MessageWithUser) error
	BroadcastNewChatName(chatID int, name string) error
	CloseConn(chatID int, connID string) error
	GetActiveUserIDs(chatID int) ([]int, error)
}

type ChatConn struct {
	UserID int
	Conn   *websocket.Conn
}

type ChatService struct {
	chats     map[int]map[string]*ChatConn
	chatsLock sync.RWMutex
}

func (s *ChatService) AddChatConn(chatID, userID int, conn *websocket.Conn) string {
	s.chatsLock.Lock()
	_, chatFound := s.chats[chatID]
	defer s.chatsLock.Unlock()
	connID := uuid.New().String()
	connObj := &ChatConn{
		UserID: userID,
		Conn:   conn,
	}
	if !chatFound {
		s.chats[chatID] = map[string]*ChatConn{
			connID: connObj,
		}
	} else {
		s.chats[chatID][connID] = connObj
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
			err := conn.Conn.Close()
			if err != nil {
				return err
			}
			delete(chatConns, connID)
		}
	}
	return ChatNotFoundErr
}

func (s *ChatService) broadcastMessage(chatID int, message any) error {
	s.chatsLock.Lock()
	defer s.chatsLock.Unlock()
	chatConns, chatFound := s.chats[chatID]
	if !chatFound {
		return ChatNotFoundErr
	}
	conns := make([]*websocket.Conn, 0)

	for _, connObj := range chatConns {
		conns = append(conns, connObj.Conn)
	}
	return broadcast(message, conns)
}

func (s *ChatService) GetActiveUserIDs(chatID int) ([]int, error) {
	s.chatsLock.Lock()
	defer s.chatsLock.Unlock()
	memberIDs := make([]int, 0)
	chatConns, chatFound := s.chats[chatID]
	if !chatFound {
		return memberIDs, ChatNotFoundErr
	}

	for _, connObj := range chatConns {
		memberIDs = append(memberIDs, connObj.UserID)
	}

	return memberIDs, nil
}

func NewChatService() *ChatService {
	return &ChatService{
		chats:     make(map[int]map[string]*ChatConn),
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

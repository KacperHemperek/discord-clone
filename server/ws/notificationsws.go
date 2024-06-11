package ws

import (
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kacperhemperek/discord-go/models"
	"sync"
)

var (
	NoUserConns = errors.New("user has no connections")
)

type NotificationService struct {
	conns     map[int]map[string]*websocket.Conn
	connsLock sync.RWMutex
}

type NotificationServiceInterface interface {
	AddConn(userID int, conn *websocket.Conn) string
	RemoveConn(userID int, connID string) error
	SendNotification(userID int, n *models.FriendRequestNotification) error
}

func (s *NotificationService) AddConn(userID int, conn *websocket.Conn) string {
	s.connsLock.Lock()
	defer s.connsLock.Unlock()
	connID := uuid.New().String()
	if userConnectionMap, userConnectionsExist := s.conns[userID]; userConnectionsExist {
		userConnectionMap[connID] = conn
	} else {
		s.conns[userID] = map[string]*websocket.Conn{
			connID: conn,
		}
	}
	return connID
}

func (s *NotificationService) RemoveConn(userID int, connID string) error {
	s.connsLock.Lock()
	defer s.connsLock.Unlock()
	if conn, connFound := s.conns[userID][connID]; connFound {
		err := conn.Close()
		delete(s.conns[userID], connID)
		if len(s.conns[userID]) == 0 {
			delete(s.conns, userID)
		}
		return err
	}
	return NoUserConns
}

func (s *NotificationService) SendNotification(userID int, n *models.FriendRequestNotification) error {
	s.connsLock.Lock()
	defer s.connsLock.Unlock()
	if conns, userConnsFound := s.conns[userID]; userConnsFound {
		for _, conn := range conns {
			err := conn.WriteJSON(*n)
			return err
		}
		return nil
	}
	return NoUserConns
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		conns:     make(map[int]map[string]*websocket.Conn),
		connsLock: sync.RWMutex{},
	}
}

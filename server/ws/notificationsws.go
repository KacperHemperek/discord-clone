package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type NotificationService struct {
	conns map[int]*websocket.Conn
}

func (s *NotificationService) AddConn(userId int, conn *websocket.Conn) {
	s.conns[userId] = conn
}

func (s *NotificationService) RemoveConn(userId int) {
	delete(s.conns, userId)
}

func (s *NotificationService) Notify(userId int, msg string) error {
	conn, ok := s.conns[userId]
	if !ok {
		return fmt.Errorf("no connection for user %d", userId)
	}

	return conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		conns: make(map[int]*websocket.Conn),
	}
}

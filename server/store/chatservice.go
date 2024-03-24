package store

import (
	"fmt"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/utils"
	"strings"
	"time"
)

type ChatService struct {
	db *Database
}

type ChatServiceInterface interface {
	GetPrivateChatByUserIDs(int, int) (*models.Chat, error)
	CreatePrivateChatWithUsers(int, int) (*models.Chat, error)
	GetUsersChatsWithMembers(userID int) ([]*models.ChatWithMembers, error)
	CreateGroupChat(chatName string, userIDs []int) (*models.Chat, error)
}

func (s *ChatService) GetPrivateChatByUserIDs(userOneID, userTwoID int) (*models.Chat, error) {
	defer utils.LogServiceCall("ChatService", "GetPrivateChatByUserIDs", time.Now())
	row := s.db.QueryRow(
		`SELECT chats.id, chats.name, chats.type, chats.created_at, chats.updated_at
			FROM chats
    		JOIN public.chat_to_user ctu1 on chats.id = ctu1.chat_id AND ctu1.user_id = $1
    		JOIN public.chat_to_user ctu2 on chats.id = ctu1.chat_id AND ctu2.user_id = $2`,
		userOneID,
		userTwoID,
	)

	chat, err := scanChat(row)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) CreatePrivateChatWithUsers(userOneID, userTwoID int) (*models.Chat, error) {
	defer utils.LogServiceCall("ChatService", "CreatePrivateChatWithUsers", time.Now())
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(
		"INSERT INTO chats (type, name) VALUES('private', 'privchat') RETURNING id, name, type, created_at, updated_at;",
	)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			fmt.Println("Error when rolling back changes: ", err)
		}
		return nil, err
	}
	var chat *models.Chat
	for rows.Next() {
		chat, err = scanChat(rows)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				fmt.Println("Error when rolling back changes: ", err)
			}
			return nil, err
		}
	}
	_, err = tx.Query(
		"INSERT INTO chat_to_user (chat_id, user_id) VALUES($1, $2),($1, $3)",
		chat.ID,
		userOneID,
		userTwoID,
	)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			fmt.Println("Error when rolling back changes: ", err)
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *ChatService) GetUsersChatsWithMembers(userID int) ([]*models.ChatWithMembers, error) {
	rows, err := s.db.Query(
		"SELECT chats.id, chats.name, chats.type, chats.created_at, chats.updated_at FROM chats JOIN chat_to_user member on user_id=$1 WHERE chats.id = member.chat_id;",
		userID,
	)
	if err != nil {
		return make([]*models.ChatWithMembers, 0), nil
	}
	chats := make([]*models.Chat, 0)
	for rows.Next() {
		chat, err := scanChat(rows)

		if err != nil {
			return make([]*models.ChatWithMembers, 0), err
		}
		chats = append(chats, chat)
	}
	result := make([]*models.ChatWithMembers, 0)
	for _, chat := range chats {
		members := make([]*models.User, 0)
		rows, err := s.db.Query(
			"SELECT users.id, users.username, users.email, users.active, users.password, users.created_at, users.updated_at FROM users JOIN chat_to_user member on member.chat_id = $1 WHERE users.id = member.user_id",
			chat.ID,
		)
		if err != nil {
			return make([]*models.ChatWithMembers, 0), err
		}
		for rows.Next() {
			member, err := scanUser(rows)
			if err != nil {
				return make([]*models.ChatWithMembers, 0), err
			}
			members = append(members, member)
		}
		chatWithMembers := &models.ChatWithMembers{
			Members: members,
			Chat:    *chat,
		}

		result = append(result, chatWithMembers)
	}
	return result, nil
}

func (s *ChatService) CreateGroupChat(chatName string, userIDs []int) (*models.Chat, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	row := tx.QueryRow(
		"INSERT INTO chats (name, type) VALUES ($1, 'group') RETURNING id, name, type,created_at, updated_at",
		chatName,
	)
	chat, err := scanChat(row)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			fmt.Println("Error rolling back transaction")
		}
		return nil, err
	}
	userVals := make([]string, len(userIDs))
	for i, userID := range userIDs {
		userVals[i] = fmt.Sprintf("(%d, %d)", chat.ID, userID)
	}
	query := fmt.Sprintf("INSERT INTO chat_to_user (chat_id, user_id) VALUES %s", strings.Join(userVals, ","))
	_, err = tx.Exec(query)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			fmt.Println("Error rolling back transaction")
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return chat, err
}
func NewChatService(db *Database) *ChatService {
	return &ChatService{
		db: db,
	}
}

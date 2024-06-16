package store

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
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
	GetChatByID(chatID int) (*models.Chat, error)
	EnrichChatWithMessages(chat *models.Chat) (*models.ChatWithMessages, error)
	GetChatMembersExcluding(chatID int, excludeUserIDs []int) ([]*models.User, error)
	GetChatMembers(chatID int) ([]*models.User, error)
	UpdateChatName(chatID int, newName string) error
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
	tx, err := s.db.Begin()
	defer func(now time.Time) {
		utils.LogServiceCall("ChatService", "GetUsersChatsWithMembers", now)
		rollback(tx)
	}(time.Now())

	rows, err := s.db.Query(
		"SELECT chats.id, chats.name, chats.type, chats.created_at, chats.updated_at FROM chats JOIN chat_to_user member on user_id=$1 WHERE chats.id = member.chat_id ORDER BY chats.updated_at DESC;",
		userID,
	)
	if err != nil {

		return make([]*models.ChatWithMembers, 0), nil
	}

	chats := make([]*models.ChatWithMembers, 0)

	for rows.Next() {
		chat, err := scanChat(rows)

		if err != nil {
			return make([]*models.ChatWithMembers, 0), err
		}

		members, err := s.getChatMembers(tx, &GetMembersFilters{
			ExcludedIDs: make([]int, 0),
			ChatID:      chat.ID,
		})
		if err != nil {
			return make([]*models.ChatWithMembers, 0), err
		}

		chatWithMembers := &models.ChatWithMembers{
			Members: members,
			Chat:    *chat,
		}

		chats = append(chats, chatWithMembers)

	}
	return chats, nil
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

func (s *ChatService) GetChatByID(chatID int) (*models.Chat, error) {
	row := s.db.QueryRow(
		"SELECT id, name, type,created_at, updated_at FROM chats WHERE id = $1",
		chatID,
	)
	return scanChat(row)
}

func (s *ChatService) EnrichChatWithMessages(chat *models.Chat) (*models.ChatWithMessages, error) {
	rows, err := s.db.Query(`
		SELECT m.id,
		       m.text, 
		       m.created_at,
		       m.updated_at, 
		       u.id, 
		       u.username,
		       u.email,
		       u.active,
		       u.password,
		       u.created_at, 
		       u.updated_at 
		FROM messages m JOIN users u on u.id = m.sender_id  WHERE m.chat_id = $1 
		ORDER BY m.created_at DESC`,
		chat.ID,
	)
	if err != nil {
		return nil, err
	}
	messages := make([]*models.MessageWithUser, 0)
	for rows.Next() {
		m, err := scanMessageWithUser(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	cwm := &models.ChatWithMessages{
		Messages: messages,
		Chat:     *chat,
	}
	return cwm, nil
}

func (s *ChatService) GetChatMembers(chatID int) ([]*models.User, error) {
	tx, err := s.db.Begin()

	defer func(now time.Time) {
		utils.LogServiceCall("ChatService", "GetChatMembers", now)
		rollback(tx)
	}(time.Now())

	if err != nil {
		return make([]*models.User, 0), err
	}
	return s.getChatMembers(tx, &GetMembersFilters{
		ExcludedIDs: make([]int, 0),
		ChatID:      chatID,
	})
}

func (s *ChatService) GetChatMembersExcluding(chatID int, excludeUserIDs []int) ([]*models.User, error) {
	tx, err := s.db.Begin()

	defer func(now time.Time) {
		utils.LogServiceCall("ChatService", "GetChatMembersExcluding", now)
		rollback(tx)
	}(time.Now())
	if err != nil {
		return make([]*models.User, 0), err
	}
	members, err := s.getChatMembers(tx, &GetMembersFilters{
		ExcludedIDs: excludeUserIDs,
		ChatID:      chatID,
	})

	if err != nil {
		return make([]*models.User, 0), err
	}

	err = tx.Commit()

	if err != nil {
		return make([]*models.User, 0), err
	}

	return members, nil
}

func (s *ChatService) UpdateChatName(chatID int, newName string) error {
	_, err := s.db.Exec(
		"UPDATE chats SET name = $1, updated_at=current_timestamp WHERE id = $2",
		newName,
		chatID,
	)
	return err
}

func (s *ChatService) getChats(tx *sql.Tx, filter *GetChatsFilters) ([]*models.Chat, error) {
	return nil, nil
}

func (s *ChatService) getChatMembers(tx *sql.Tx, filter *GetMembersFilters) ([]*models.User, error) {
	where := []string{
		"cu.chat_id = @chat_id",
	}
	args := pgx.NamedArgs{
		"chat_id": filter.ChatID,
	}

	if v := filter.ExcludedIDs; len(v) != 0 {
		argList := make([]string, 0)
		for i, ID := range v {
			argID := fmt.Sprintf("@user_id_%d", i)
			args[argID[1:]] = ID
			argList = append(argList, argID)
		}
		idList := "(" + strings.Join(argList, ",") + ")"

		where = append(where, "cu.user_id NOT IN "+idList)
	}

	rows, err := tx.Query(`
		SELECT 
			u.id, u.username, u.email, u.active, u.password, u.created_at, u.updated_at 
		FROM 
			chat_to_user cu
		JOIN users u on u.id = cu.user_id `+
		whereSQL(where)+";",
		args,
	)

	members := make([]*models.User, 0)

	if err != nil {
		return members, nil
	}

	for rows.Next() {
		member, err := scanUser(rows)
		if err != nil {
			return make([]*models.User, 0), err
		}
		members = append(members, member)
	}

	return members, nil
}

type GetMembersFilters struct {
	ExcludedIDs []int
	ChatID      int
}

type GetChatsFilters struct {
	ChatID *int
}

func NewChatService(db *Database) *ChatService {
	return &ChatService{
		db: db,
	}
}

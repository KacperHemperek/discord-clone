package store

import (
	"github.com/jackc/pgx/v5"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/utils"
	"time"
)

type MessageServiceInterface interface {
	CreateMessageInChat(chatID, userID int, text string) (*models.Message, error)
	EnrichMessageWithUser(message *models.Message) (*models.MessageWithUser, error)
}

type MessageService struct {
	db *Database
}

func (s *MessageService) CreateMessageInChat(chatID, userID int, text string) (*models.Message, error) {
	tx, err := s.db.Begin()
	defer func(now time.Time) {
		utils.LogServiceCall("MessageService", "CreateMessageInChat", now)
		rollback(tx)
	}(time.Now())

	if err != nil {
		return nil, err
	}

	row := tx.QueryRow(`
		INSERT INTO messages (text, sender_id, chat_id) 
			VALUES ($1, $2, $3) RETURNING id, text, created_at, updated_at;`,
		text,
		userID,
		chatID,
	)

	m, err := scanMessage(row)

	_, err = tx.Exec(
		"UPDATE chats SET updated_at = now() "+whereSQL([]string{"id = @chat_id"}), pgx.NamedArgs{
			"chat_id": chatID,
		})

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *MessageService) EnrichMessageWithUser(message *models.Message) (*models.MessageWithUser, error) {
	row := s.db.QueryRow(`
		SELECT u.id, u.username, u.email, u.active, u.password, u.created_at, u.updated_at 
			FROM messages m JOIN users u on u.id = m.sender_id WHERE m.id = $1`,
		message.ID,
	)
	user, err := scanUser(row)
	if err != nil {
		return nil, err
	}
	mwu := &models.MessageWithUser{
		User:    user,
		Message: *message,
	}
	return mwu, nil
}

func NewMessageService(db *Database) *MessageService {
	return &MessageService{
		db: db,
	}
}

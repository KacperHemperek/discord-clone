package store

import "github.com/kacperhemperek/discord-go/models"

type MessageServiceInterface interface {
	CreateMessageInChat(chatID, userID int, text string) (*models.Message, error)
	EnrichMessageWithUser(chat *models.Message) (*models.MessageWithUser, error)
}

type MessageService struct {
	db *Database
}

func (s *MessageService) CreateMessageInChat(chatID, userID int, text string) (*models.Message, error) {
	row := s.db.QueryRow(`
		INSERT INTO messages (text, sender_id, chat_id) 
			VALUES ($1, $2, $3) RETURNING id, text, created_at, updated_at;`,
		text,
		userID,
		chatID,
	)

	return scanMessage(row)
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

package store

import (
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/types"
)

type NotificationServiceInterface interface {
	CreateFriendRequestNotification(userID int, data models.FriendRequestNotificationData) (*models.FriendRequestNotification, error)
	GetUserFriendRequestNotifications(userID int, seen *string) ([]*models.FriendRequestNotification, error)
}

type NotificationService struct {
	db        *Database
	validator *validator.Validate
}

func (s *NotificationService) CreateFriendRequestNotification(userID int, data models.FriendRequestNotificationData) (*models.FriendRequestNotification, error) {

	jsonData, jsonMarshalError := json.Marshal(data)

	if jsonMarshalError != nil {
		return nil, jsonMarshalError
	}

	row := s.db.QueryRow(
		"INSERT INTO notifications (user_id, type, data) VALUES ($1, $2, $3) RETURNING id, type, user_id, data, seen, created_at, updated_at;",
		userID,
		types.FriendRequestNotification.String(),
		jsonData,
	)
	notificationDto := &models.NotificationDTO{}
	err := row.Scan(
		&notificationDto.BaseNotification.Base.ID,
		&notificationDto.BaseNotification.Type,
		&notificationDto.BaseNotification.UserID,
		&notificationDto.Data,
		&notificationDto.BaseNotification.Seen,
		&notificationDto.BaseNotification.Base.CreatedAt,
		&notificationDto.BaseNotification.Base.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	notificationData := &models.FriendRequestNotificationData{}

	if err := json.Unmarshal(notificationDto.Data, notificationData); err != nil {
		return nil, err
	}

	if dataValidationError := s.validator.Struct(notificationData); dataValidationError != nil {
		return nil, dataValidationError
	}

	return &models.FriendRequestNotification{
		BaseNotification: notificationDto.BaseNotification,
		Data:             notificationData,
	}, nil
}

func (s *NotificationService) GetUserFriendRequestNotifications(userID int, seen *string) ([]*models.FriendRequestNotification, error) {
	var err error
	tx, err := s.db.Begin()
	defer func() {
		if err != nil {
			rollback(tx)
		}
		tx.Commit()
	}()
	return s.findFriendRequestNotifications(tx, &NotificationFilters{
		Seen:   seen,
		UserID: &userID,
	})
}

func (s *NotificationService) findFriendRequestNotifications(tx *sql.Tx, filters *NotificationFilters) ([]*models.FriendRequestNotification, error) {
	where := []string{
		"type = @type",
	}
	args := pgx.NamedArgs{
		"type": types.FriendRequestNotification.String(),
	}

	if v := filters.Seen; v != nil {
		where = append(where, "seen = @seen")
		args["seen"] = v
	}

	if v := filters.UserID; v != nil {
		where = append(where, "user_id = @user_id")
		args["user_id"] = v
	}

	rows, err := tx.Query(
		"SELECT id, type, user_id, data, seen, created_at, updated_at FROM notifications "+whereSQL(where)+" ORDER BY created_at DESC",
		args,
	)

	notifications := make([]*models.FriendRequestNotification, 0)

	if err != nil {
		return notifications, err
	}

	for rows.Next() {
		n, err := scanFriendRequestNotification(rows)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

type NotificationFilters struct {
	Seen   *string
	UserID *int
}

func NewNotificationService(db *Database, v *validator.Validate) *NotificationService {
	return &NotificationService{
		db:        db,
		validator: v,
	}
}

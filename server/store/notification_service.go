package store

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/types"
)

type NotificationServiceInterface interface {
	CreateFriendRequestNotification(userID int, data models.FriendRequestNotificationData) (*models.FriendRequestNotification, error)
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

func NewNotificationService(db *Database, v *validator.Validate) *NotificationService {
	return &NotificationService{
		db:        db,
		validator: v,
	}
}

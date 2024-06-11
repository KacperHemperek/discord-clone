package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/types"
	"log/slog"
)

type NotificationServiceInterface interface {
	CreateFriendRequestNotification(userID int, data models.FriendRequestNotificationData) (*models.FriendRequestNotification, error)
	GetUserFriendRequestNotifications(userID int, seen *BoolFilter) ([]*models.FriendRequestNotification, error)
	MarkUsersNotificationsAsSeenByType(userID int, nType string) error
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

func (s *NotificationService) GetUserFriendRequestNotifications(userID int, seen *BoolFilter) ([]*models.FriendRequestNotification, error) {
	tx, err := s.db.Begin()
	defer func() {
		if err = tx.Rollback(); err != nil && errors.Is(err, sql.ErrTxDone) {
			slog.Error("rollback tx", "error", err)
		}
	}()
	if err != nil {
		return make([]*models.FriendRequestNotification, 0), err
	}

	notifications, err := s.findFriendRequestNotifications(tx, &FindNotificationFilters{
		Seen:   seen,
		UserID: &userID,
	})
	if err != nil {
		return notifications, err
	}
	err = tx.Commit()
	if err != nil {
		slog.Error("commit tx", "error", err)
	}
	return notifications, nil
}

func (s *NotificationService) MarkUsersNotificationsAsSeenByType(userID int, nType string) error {
	tx, err := s.db.Begin()
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("rollback tx", "error", err)
		}
	}()
	if err != nil {
		return err
	}

	seen, err := NewBoolFilter("false")
	if err != nil {
		return err
	}

	err = s.markAsSeen(tx, &UpdateNotificationFilters{
		UserID: &userID,
		Type:   &nType,
		Seen:   seen,
	})
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *NotificationService) markAsSeen(tx *sql.Tx, filters *UpdateNotificationFilters) error {

	where := make([]string, 0)
	args := pgx.NamedArgs{}

	if v := filters.Type; v != nil && types.IsNotificationType(*v) {
		where = append(where, "type = @type")
		args["type"] = v
	}

	if v := filters.UserID; v != nil {
		where = append(where, "user_id = @user_id")
		args["user_id"] = v
	}

	if v := filters.Seen; v != nil {
		where = append(where, "seen = @seen")
		args["seen"] = v
	}

	_, err := tx.Exec(
		"UPDATE notifications SET seen = true "+
			whereSQL(where)+";",
		args,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *NotificationService) findFriendRequestNotifications(tx *sql.Tx, filters *FindNotificationFilters) ([]*models.FriendRequestNotification, error) {
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
		"SELECT id, type, user_id, data, seen, created_at, updated_at FROM notifications "+
			whereSQL(where)+
			" ORDER BY created_at DESC",
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

type FindNotificationFilters struct {
	Seen   *BoolFilter
	UserID *int
}

type UpdateNotificationFilters struct {
	UserID *int
	Type   *string
	Seen   *BoolFilter
}

func NewNotificationService(db *Database, v *validator.Validate) *NotificationService {
	return &NotificationService{
		db:        db,
		validator: v,
	}
}

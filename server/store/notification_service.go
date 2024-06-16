package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/types"
	"github.com/kacperhemperek/discord-go/utils"
	"log/slog"
	"strings"
	"time"
)

type NotificationServiceInterface interface {
	GetUserFriendRequestNotifications(userID int, seen *BoolFilter, limit *LimitFilter) ([]*models.FriendRequestNotification, error)
	GetUserNewMessageNotifications(userID int, seen *BoolFilter, limit *LimitFilter) ([]*models.NewMessageNotification, error)

	CreateFriendRequestNotification(userID int, data models.FriendRequestNotificationData) (*models.FriendRequestNotification, error)
	CreateNewMessageNotificationsForUsers(userIDs []int, data *models.NewMessageNotificationData) ([]*models.NewMessageNotification, error)

	MarkUsersNotificationsAsSeen(userID int, nType string) error
	MarkUsersNewMessageNotificationsAsSeenByChatID(userID, chatID int) error
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

	frn := types.FriendRequestNotification
	row := s.db.QueryRow(
		"INSERT INTO notifications (user_id, type, data) VALUES ($1, $2, $3) RETURNING id, type, user_id, data, seen, created_at, updated_at;",
		userID,
		frn.String(),
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

func (s *NotificationService) GetUserFriendRequestNotifications(
	userID int,
	seen *BoolFilter,
	limit *int,
) ([]*models.FriendRequestNotification, error) {

	tx, err := s.db.Begin()

	defer func(now time.Time) {
		utils.LogServiceCall("NotificationsService", "GetUserFriendRequestNotifications", now)
		rollback(tx)
	}(time.Now())

	if err != nil {
		return make([]*models.FriendRequestNotification, 0), err
	}

	notifications, err := s.findFriendRequestNotifications(tx, &FindNotificationFilters{
		Seen:   seen,
		UserID: &userID,
		Limit:  limit,
	})

	if err != nil {
		return notifications, err
	}
	err = tx.Commit()
	if err != nil {
		return make([]*models.FriendRequestNotification, 0), err
	}
	return notifications, nil
}

func (s *NotificationService) MarkUsersNotificationsAsSeen(userID int, nType string) error {
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

func (s *NotificationService) CreateNewMessageNotificationsForUsers(
	userIDs []int,
	data *models.NewMessageNotificationData,
) ([]*models.NewMessageNotification, error) {
	defer func(now time.Time) {
		utils.LogServiceCall("NotificationService", "CreateNewMessageNotificationsForUsers", now)
	}(time.Now())

	jsonData, jsonMarshalError := json.Marshal(data)

	if jsonMarshalError != nil {
		return nil, jsonMarshalError
	}

	values := make([]string, 0)
	nmn := types.NewMessageNotification

	for _, userID := range userIDs {
		values = append(
			values,
			fmt.Sprintf("(%d, false, @json_data, @type)", userID),
		)
	}

	if len(values) == 0 {
		slog.Info("no notifications to create")
		return make([]*models.NewMessageNotification, 0), nil
	}

	rows, err := s.db.Query(
		"INSERT INTO notifications (user_id, seen, data, type) VALUES "+strings.Join(values, ",")+" RETURNING id, type, seen, data, user_id, created_at, updated_at;",
		pgx.NamedArgs{
			"json_data": jsonData,
			"type":      nmn.String(),
		},
	)

	ns := make([]*models.NewMessageNotification, 0)

	if err != nil {
		return ns, err
	}

	for rows.Next() {
		n := &models.NewMessageNotification{}
		err := rows.Scan(
			&n.ID,
			&n.Type,
			&n.Seen,
			&n.Data,
			&n.UserID,
			&n.CreatedAt,
			&n.UpdatedAt,
		)
		if err != nil {
			return make([]*models.NewMessageNotification, 0), err
		}
		ns = append(ns, n)
	}

	return ns, nil
}

func (s *NotificationService) GetUserNewMessageNotifications(userID int, seen *BoolFilter, limit *LimitFilter) ([]*models.NewMessageNotification, error) {
	tx, err := s.db.Begin()

	defer func(now time.Time) {
		utils.LogServiceCall("NotificationsService", "GetUserNewMessageNotifications", now)
		rollback(tx)
	}(time.Now())

	if err != nil {
		return make([]*models.NewMessageNotification, 0), err
	}

	ns, err := s.findNewMessageNotifications(tx, &FindNotificationFilters{
		Seen:   seen,
		UserID: &userID,
		Limit:  limit,
	})

	if err != nil {
		return make([]*models.NewMessageNotification, 0), err
	}

	if err := tx.Commit(); err != nil {
		return make([]*models.NewMessageNotification, 0), err
	}

	return ns, nil
}

func (s *NotificationService) MarkUsersNewMessageNotificationsAsSeenByChatID(userID, chatID int) error {
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

	nmn := types.NewMessageNotification
	nmnStr := nmn.String()
	userNewMessageNotifications, err := s.findNewMessageNotifications(tx, &FindNotificationFilters{
		Seen:   seen,
		UserID: &userID,
	})
	if err != nil {
		return err
	}

	for _, n := range userNewMessageNotifications {
		if n.Data.ChatID != chatID {
			continue
		}
		err = s.markAsSeen(tx, &UpdateNotificationFilters{
			UserID:         &userID,
			Type:           &nmnStr,
			Seen:           seen,
			NotificationID: &n.ID,
		})
		if err != nil {
			return err
		}
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

	if v := filters.NotificationID; v != nil {
		where = append(where, "id = @id")
		args["id"] = v
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
	frn := types.FriendRequestNotification
	where := []string{
		"type = @type",
	}
	limit := ""
	args := pgx.NamedArgs{
		"type": frn.String(),
	}

	if v := filters.Seen; v != nil {
		where = append(where, "seen = @seen")
		args["seen"] = v
	}

	if v := filters.UserID; v != nil {
		where = append(where, "user_id = @user_id")
		args["user_id"] = v
	}

	if v := filters.Limit; v != nil {
		limit = fmt.Sprintf(" LIMIT %d", *v)
	}

	rows, err := tx.Query(
		"SELECT id, type, user_id, data, seen, created_at, updated_at FROM notifications "+
			whereSQL(where)+
			" ORDER BY created_at DESC"+
			limit+
			";",
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

func (s *NotificationService) findNewMessageNotifications(tx *sql.Tx, filters *FindNotificationFilters) ([]*models.NewMessageNotification, error) {
	nmn := types.NewMessageNotification
	where := []string{
		"type = @type",
	}
	limit := ""
	args := pgx.NamedArgs{
		"type": nmn.String(),
	}

	if v := filters.Seen; v != nil {
		where = append(where, "seen = @seen")
		args["seen"] = v
	}

	if v := filters.UserID; v != nil {
		where = append(where, "user_id = @user_id")
		args["user_id"] = v
	}

	if v := filters.Limit; v != nil {
		limit = fmt.Sprintf(" LIMIT %d", *v)
	}

	rows, err := tx.Query(
		"SELECT id, type, user_id, data, seen, created_at, updated_at FROM notifications "+
			whereSQL(where)+
			" ORDER BY created_at DESC"+
			limit+
			";",
		args,
	)

	notifications := make([]*models.NewMessageNotification, 0)
	if err != nil {
		return notifications, err
	}
	for rows.Next() {
		n := &models.NewMessageNotification{}
		err := rows.Scan(
			&n.ID,
			&n.Type,
			&n.UserID,
			&n.Data,
			&n.Seen,
			&n.CreatedAt,
			&n.UpdatedAt,
		)
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
	Limit  *LimitFilter
}

type UpdateNotificationFilters struct {
	UserID         *int
	Type           *string
	Seen           *BoolFilter
	NotificationID *int
}

func NewNotificationService(db *Database, v *validator.Validate) *NotificationService {
	return &NotificationService{
		db:        db,
		validator: v,
	}
}

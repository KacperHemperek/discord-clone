package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/utils"
	"time"
)

type FriendshipService struct {
	db *Database
}

type FriendshipServiceInterface interface {
	SendFriendRequest(inviterId, friendId int) error
	GetUsersFriendRequests(userId int) ([]*models.FriendRequest, error)
	GetFriendshipByUsers(userOneID, userTwoID int) (*models.Friendship, error)
	GetFriendshipById(requestId int) (*models.Friendship, error)
	AcceptFriendRequest(requestId int) error
	RejectFriendRequest(requestId int) error
	MakeFriendshipPending(requestId int) error
	DeleteRequestAndSendNew(requestId, inviterId, friendId int) error
}

func (s *FriendshipService) SendFriendRequest(inviterId, friendId int) error {
	defer utils.LogServiceCall("FriendshipService", "SendFriendRequest", time.Now())
	ctx := context.Background()
	context.WithTimeout(ctx, 200*time.Millisecond)
	_, err := s.db.QueryContext(
		ctx,
		"INSERT INTO friendships (inviter_id, friend_id) VALUES ($1, $2)",
		inviterId, friendId)

	if err != nil {
		return err
	}

	return nil
}

func (s *FriendshipService) GetUsersFriendRequests(userId int) ([]*models.FriendRequest, error) {
	defer utils.LogServiceCall("FriendshipService", "GetUsersFriendRequests", time.Now())
	ctx := context.Background()
	context.WithTimeout(ctx, 200*time.Millisecond)
	rows, err := s.db.QueryContext(
		ctx,
		`
		SELECT f.id, f.status, f.requested_at, f.status_updated_at, u.id, u.username, u.email, u.active, u.password, u.created_at, u.updated_at 
		FROM friendships f JOIN users u ON f.inviter_id = u.id WHERE f.friend_id = $1 AND f.status = 'pending';
		`,
		userId,
	)
	users := make([]*models.FriendRequest, 0)
	if err != nil {
		return users, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			fmt.Println("ERROR userservice: ", err)
		}
	}()

	for rows.Next() {
		user, err := scanFriendRequest(rows)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return []*models.FriendRequest{}, nil
			}
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *FriendshipService) GetFriendshipByUsers(userOneID, userTwoID int) (*models.Friendship, error) {
	defer utils.LogServiceCall("FriendshipService", "GetFriendshipByUsers", time.Now())
	ctx := context.Background()
	context.WithTimeout(ctx, 200*time.Millisecond)
	row := s.db.QueryRowContext(
		ctx,
		"SELECT id, inviter_id, friend_id, status, seen, requested_at, status_updated_at FROM friendships WHERE (inviter_id = $1 AND friend_id = $2) OR (inviter_id = $2 AND friend_id = $1);",
		userOneID, userTwoID,
	)

	friendship, err := scanFriendship(row)

	if err != nil {
		return nil, err
	}

	return friendship, nil
}

func (s *FriendshipService) GetFriendshipById(requestId int) (*models.Friendship, error) {
	defer utils.LogServiceCall("FriendshipService", "GetFriendshipById", time.Now())
	ctx := context.Background()
	context.WithTimeout(ctx, 200*time.Millisecond)
	row := s.db.QueryRowContext(
		ctx,
		"SELECT id, inviter_id, friend_id, status, seen, requested_at, status_updated_at FROM friendships WHERE id = $1;",
		requestId,
	)

	friendship, err := scanFriendship(row)

	if err != nil {
		return nil, err
	}

	return friendship, nil
}

func (s *FriendshipService) AcceptFriendRequest(requestId int) error {
	defer utils.LogServiceCall("FriendshipService", "AcceptFriendRequest", time.Now())
	ctx := context.Background()
	context.WithTimeout(ctx, 200*time.Millisecond)
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE friendships SET status = 'accepted', status_updated_at = CURRENT_TIMESTAMP WHERE id = $1;",
		requestId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *FriendshipService) RejectFriendRequest(requestId int) error {
	defer utils.LogServiceCall("FriendshipService", "RejectFriendRequest", time.Now())
	ctx := context.Background()
	context.WithTimeout(ctx, 200*time.Millisecond)
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE friendships SET status = 'rejected', status_updated_at = CURRENT_TIMESTAMP WHERE id = $1;",
		requestId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *FriendshipService) MakeFriendshipPending(requestId int) error {
	defer utils.LogServiceCall("FriendshipService", "MakeFriendshipPending", time.Now())
	ctx := context.Background()
	context.WithTimeout(ctx, 200*time.Millisecond)

	_, err := s.db.ExecContext(
		ctx,
		"UPDATE friendships SET status = 'pending', status_updated_at = CURRENT_TIMESTAMP WHERE id = $1;",
		requestId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *FriendshipService) DeleteRequestAndSendNew(requestId, inviterId, friendId int) error {
	defer utils.LogServiceCall("FriendshipService", "DeleteRequestAndSendNew", time.Now())
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"DELETE FROM friendships WHERE id = $1;",
		requestId,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO friendships (inviter_id, friend_id) VALUES ($1, $2);",
		inviterId, friendId,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func NewFriendshipService(db *Database) *FriendshipService {
	return &FriendshipService{db: db}
}

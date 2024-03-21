package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/utils"
	"sort"
	"sync"
	"time"
)

type FriendshipService struct {
	db *Database
}

type FriendshipServiceInterface interface {
	SendFriendRequest(inviterID, friendID int) error
	GetUsersFriendRequests(userID int) ([]*models.FriendRequest, error)
	GetFriendshipByUsers(userOneID, userTwoID int) (*models.Friendship, error)
	GetFriendshipByID(requestID int) (*models.Friendship, error)
	AcceptFriendRequest(requestID int) error
	RejectFriendRequest(requestID int) error
	MakeFriendshipPending(requestID int) error
	DeleteRequestAndSendNew(requestID, inviterID, friendID int) error
	GetFriendsByUserID(userID int) ([]*models.User, error)
}

func (s *FriendshipService) SendFriendRequest(inviterID, friendID int) error {
	defer utils.LogServiceCall("FriendshipService", "SendFriendRequest", time.Now())
	ctx := context.Background()
	context.WithTimeout(ctx, 200*time.Millisecond)
	_, err := s.db.QueryContext(
		ctx,
		"INSERT INTO friendships (inviter_id, friend_id) VALUES ($1, $2)",
		inviterID, friendID)

	if err != nil {
		return err
	}

	return nil
}

func (s *FriendshipService) GetUsersFriendRequests(userID int) ([]*models.FriendRequest, error) {
	defer utils.LogServiceCall("FriendshipService", "GetUsersFriendRequests", time.Now())
	ctx := context.Background()
	context.WithTimeout(ctx, 200*time.Millisecond)
	rows, err := s.db.QueryContext(
		ctx,
		`
		SELECT f.id, f.status, f.requested_at, f.status_updated_at, u.id, u.username, u.email, u.active, u.password, u.created_at, u.updated_at 
		FROM friendships f JOIN users u ON f.inviter_id = u.id WHERE f.friend_id = $1 AND f.status = 'pending';
		`,
		userID,
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

func (s *FriendshipService) GetFriendshipByID(requestID int) (*models.Friendship, error) {
	defer utils.LogServiceCall("FriendshipService", "GetFriendshipByID", time.Now())
	ctx := context.Background()
	context.WithTimeout(ctx, 200*time.Millisecond)
	row := s.db.QueryRowContext(
		ctx,
		"SELECT id, inviter_id, friend_id, status, seen, requested_at, status_updated_at FROM friendships WHERE id = $1;",
		requestID,
	)

	friendship, err := scanFriendship(row)

	if err != nil {
		return nil, err
	}

	return friendship, nil
}

func (s *FriendshipService) AcceptFriendRequest(requestID int) error {
	defer utils.LogServiceCall("FriendshipService", "AcceptFriendRequest", time.Now())
	ctx := context.Background()
	context.WithTimeout(ctx, 200*time.Millisecond)
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE friendships SET status = 'accepted', status_updated_at = CURRENT_TIMESTAMP WHERE id = $1;",
		requestID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *FriendshipService) RejectFriendRequest(requestID int) error {
	defer utils.LogServiceCall("FriendshipService", "RejectFriendRequest", time.Now())
	ctx := context.Background()
	context.WithTimeout(ctx, 200*time.Millisecond)
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE friendships SET status = 'rejected', status_updated_at = CURRENT_TIMESTAMP WHERE id = $1;",
		requestID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *FriendshipService) MakeFriendshipPending(requestID int) error {
	defer utils.LogServiceCall("FriendshipService", "MakeFriendshipPending", time.Now())
	ctx := context.Background()
	context.WithTimeout(ctx, 200*time.Millisecond)

	_, err := s.db.ExecContext(
		ctx,
		"UPDATE friendships SET status = 'pending', status_updated_at = CURRENT_TIMESTAMP WHERE id = $1;",
		requestID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *FriendshipService) DeleteRequestAndSendNew(requestID, inviterID, friendID int) error {
	defer utils.LogServiceCall("FriendshipService", "DeleteRequestAndSendNew", time.Now())
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"DELETE FROM friendships WHERE id = $1;",
		requestID,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO friendships (inviter_id, friend_id) VALUES ($1, $2);",
		inviterID, friendID,
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

func (s *FriendshipService) GetFriendsByUserID(userID int) ([]*models.User, error) {
	defer utils.LogServiceCall("FriendshipService", "GetFriendsByUserID", time.Now())

	errChan := make(chan error)
	doneChan := make(chan bool)
	friendChan := make(chan *models.Friend)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		wg.Wait()
		doneChan <- true
	}()

	go func() {
		defer wg.Done()
		acceptedFriends, err := s.db.Query(
			"SELECT u.id, u.username, u.email, u.active, u.password, u.created_at, u.updated_at, f.status_updated_at FROM friendships f JOIN public.users u on u.id = f.inviter_id WHERE friend_id=$1;",
			userID,
		)
		if err != nil {
			errChan <- err
			return
		}
		for acceptedFriends.Next() {
			friend, err := scanFriend(acceptedFriends)

			if err != nil {
				errChan <- err
				return
			}
			friendChan <- friend
		}
	}()

	go func() {
		defer wg.Done()
		invitedFriends, err := s.db.Query(
			"SELECT u.id, u.username, u.email, u.active, u.password, u.created_at, u.updated_at, f.status_updated_at FROM friendships f JOIN public.users u on u.id = f.friend_id WHERE inviter_id=$1;",
			userID,
		)
		if err != nil {
			errChan <- err
			return
		}
		for invitedFriends.Next() {
			friend, err := scanFriend(invitedFriends)
			if err != nil {
				errChan <- err
				return
			}
			friendChan <- friend
		}
	}()

	friends := make([]*models.Friend, 0)
	users := make([]*models.User, 0)

	for {
		select {
		case err := <-errChan:
			return users, err
		case friend := <-friendChan:
			friends = append(friends, friend)
		case <-doneChan:
			sort.Slice(friends, func(i, j int) bool {
				return friends[i].AcceptedAt.Before(friends[j].AcceptedAt)
			})
			for _, friend := range friends {
				user := &friend.User

				users = append(users, user)
			}
			return users, nil
		}
	}
}

func NewFriendshipService(db *Database) *FriendshipService {
	return &FriendshipService{db: db}
}

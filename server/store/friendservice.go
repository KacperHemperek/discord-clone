package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kacperhemperek/discord-go/models"
)

type FriendshipService struct {
	db *Database
}

func (s *FriendshipService) SendFriendRequest(inviterId, friendId int) error {
	_, err := s.db.Query(
		"INSERT INTO friendships (inviter_id, friend_id) VALUES ($1, $2)",
		inviterId, friendId)

	if err != nil {
		return err
	}

	return nil
}

func (s *FriendshipService) GetUsersFriendRequests(userId int) ([]*models.FriendRequest, error) {
	rows, err := s.db.Query(
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
	row := s.db.QueryRow(
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
	row := s.db.QueryRow(
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
	_, err := s.db.Exec(
		"UPDATE friendships SET status = 'accepted', status_updated_at = CURRENT_TIMESTAMP WHERE id = $1;",
		requestId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *FriendshipService) RejectFriendRequest(requestId int) error {
	_, err := s.db.Exec(
		"UPDATE friendships SET status = 'rejected', status_updated_at = CURRENT_TIMESTAMP WHERE id = $1;",
		requestId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *FriendshipService) MakeFriendshipPending(requestId int) error {
	_, err := s.db.Exec(
		"UPDATE friendships SET status = 'pending', status_updated_at = CURRENT_TIMESTAMP WHERE id = $1;",
		requestId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *FriendshipService) DeleteFriendship(requestId int) error {
	_, err := s.db.Exec(
		"DELETE FROM friendships WHERE id = $1;",
		requestId,
	)

	if err != nil {
		return err
	}

	return nil
}

func NewFriendshipService(db *Database) *FriendshipService {
	return &FriendshipService{db: db}
}

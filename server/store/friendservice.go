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

func (s *FriendshipService) GetFriendRequests(userId int) ([]*models.User, error) {
	rows, err := s.db.Query(
		`
		SELECT u.id, u.username, u.email, u.active, u.password, u.created_at, u.updated_at 
		FROM friendships f JOIN users u ON f.inviter_id = u.id WHERE f.friend_id = $1 AND f.status = 'pending';
		`,
		userId,
	)
	users := make([]*models.User, 0)
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
		user, err := scanUser(rows)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return []*models.User{}, nil
			}
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func NewFriendshipService(db *Database) *FriendshipService {
	return &FriendshipService{db: db}
}

package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kacperhemperek/discord-go/models"
)

type UserService struct {
	db *Database
}

func (s *UserService) FindUserByEmail(email string) (*models.User, error) {
	rows, err := s.db.Query(
		"SELECT id, username, email, active, password, created_at, updated_at FROM users WHERE email = $1;",
		email,
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			fmt.Println("ERROR userservice: ", err)
		}
	}()

	for rows.Next() {
		user, err := ScanUser(rows)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, UserNotFoundError
			}
			return nil, err
		}

		if user == nil {
			return nil, UserNotFoundError
		}

		return user, nil
	}

	return nil, UserNotFoundError
}

func (s *UserService) CreateUser(username, password, email string) (*models.User, error) {
	rows, err := s.db.Query(
		"INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id, username, email, active, password, created_at, updated_at;",
		username, password, email,
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			fmt.Println("ERROR userservice: ", err)
		}
	}()

	for rows.Next() {
		user, err := ScanUser(rows)

		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, UserUnknownError
}

func (s *UserService) SendFriendRequest(inviterId, friendId int) error {
	_, err := s.db.Query(
		"INSERT INTO friendships (inviter_id, friend_id) VALUES ($1, $2)",
		inviterId, friendId)

	if err != nil {
		return err
	}

	return nil
}

func ScanUser(rows *sql.Rows) (*models.User, error) {
	user := &models.User{}

	err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Active, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

var UserNotFoundError = errors.New("user not found")

var UserUnknownError = errors.New("unknown error")

func NewUserService(db *Database) *UserService {
	return &UserService{db: db}
}

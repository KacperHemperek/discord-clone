package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/utils"
	"time"
)

type UserService struct {
	db *Database
}

func (s *UserService) FindUserByEmail(email string) (*models.User, error) {
	defer utils.LogServiceCall("UserService", "FindUserByEmail", time.Now())
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
		user, err := scanUser(rows)

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
	defer utils.LogServiceCall("UserService", "CreateUser", time.Now())
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
		user, err := scanUser(rows)

		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, UserUnknownError
}

var UserNotFoundError = errors.New("user not found")

var UserUnknownError = errors.New("unknown error")

func NewUserService(db *Database) *UserService {
	return &UserService{db: db}
}

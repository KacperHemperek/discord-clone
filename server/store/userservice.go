package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kacperhemperek/discord-go/models"
	"github.com/kacperhemperek/discord-go/utils"
	"strings"
	"time"
)

type UserService struct {
	db *Database
}

type UserServiceInterface interface {
	FindUserByEmail(email string) (*models.User, error)
	CreateUser(username, password, email string) (*models.User, error)
	GetUsersByIDs(userIDs []int) ([]*models.User, error)
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

func (s *UserService) GetUsersByIDs(userIDs []int) ([]*models.User, error) {
	placeholders := make([]string, len(userIDs))
	for i := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("SELECT id, username, email, active, password, created_at, updated_at FROM users WHERE id IN (%s)", strings.Join(placeholders, ", "))

	params := make([]any, len(userIDs))
	for i, id := range userIDs {
		params[i] = id
	}

	users := make([]*models.User, 0)
	rows, err := s.db.Query(query, params...)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return make([]*models.User, 0), err
		}
		users = append(users, user)
	}

	return users, nil
}

var UserNotFoundError = errors.New("user not found")

var UserUnknownError = errors.New("unknown error")

func NewUserService(db *Database) *UserService {
	return &UserService{db: db}
}

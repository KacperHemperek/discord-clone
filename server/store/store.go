package store

import "fmt"

type Database struct {
	Users []*User
}

var db = &Database{
	Users: make([]*User, 0, 10),
}

func GetUserByUsername(username string) (*User, error) {
	for _, user := range db.Users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user with username %s not found", username)
}

func GetUserByEmail(email string) (*User, error) {
	for _, user := range db.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user with email %s not found", email)
}

func CreateUser(username, email, password string) (*User, error) {
	user := &User{
		ID:       len(db.Users) + 1,
		Username: username,
		Email:    email,
		Password: password,
	}
	db.Users = append(db.Users, user)
	return user, nil
}

type Limiter struct {
	Ips map[string]string
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

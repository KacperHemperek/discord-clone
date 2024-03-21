package models

type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
	Base
}

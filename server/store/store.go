package store

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type Database struct {
	Users []*User
}

func Init() *gorm.DB {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	fmt.Printf("username: %s, password: %s, name: %s", username, password, name)
	connStr := fmt.Sprintf("user=%s password=%s host=localhost dbname=%s sslmode=disable", username, password, name)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %s", err.Error()))
	}

	return db
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

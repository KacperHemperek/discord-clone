package store

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
)

func NewDB() *Database {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dburl := fmt.Sprintf("user=%s password=%s dbname=%s host=localhost sslmode=disable", username, password, dbname)

	db, err := sql.Open("pgx", dburl)

	if err != nil {
		fmt.Println("Error connecting to database")
		panic(err)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Error pinging database")
		panic(err)
	}

	return db
}

type Database = sql.DB

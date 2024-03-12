package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

func RunMigrations(db *Database) {
	fmt.Println("Running migrations...")
	driver, err := pgx.WithInstance(db, &pgx.Config{})

	if err != nil {
		fmt.Println("Error creating migration driver")
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://store/migrations",
		"pgx", driver)

	if err != nil {
		fmt.Println("Could not create migration instance")
		panic(err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("Error running migration")
			panic(err)
		}
	}
	fmt.Println("All migrations ran successfully")
}

type Database = sql.DB

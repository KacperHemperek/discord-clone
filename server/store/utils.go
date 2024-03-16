package store

import (
	"database/sql"
	"github.com/kacperhemperek/discord-go/models"
)

func scanUser(rows *sql.Rows) (*models.User, error) {
	user := &models.User{}

	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Active,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

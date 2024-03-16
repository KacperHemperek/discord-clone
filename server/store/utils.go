package store

import (
	"github.com/kacperhemperek/discord-go/models"
)

// Scanner is a helper interface that wraps the Scan method implemented by sql.Rows and sql.Row
type Scanner interface {
	Scan(dest ...any) error
}

// Scans user from the given Scanner, this function accepts sql.Rows and sql.Row as Scanner
// and returns a user or an error if the scan fails
func scanUser(rows Scanner) (*models.User, error) {
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

// Scans friendship from the given Scanner, this function accepts sql.Rows and sql.Row as Scanner
// and returns a friendship or an error if the scan fails
func scanFriendship(row Scanner) (*models.Friendship, error) {
	friendship := &models.Friendship{}

	err := row.Scan(
		&friendship.ID,
		&friendship.InviterID,
		&friendship.FriendID,
		&friendship.Status,
		&friendship.Accepted,
		&friendship.Seen,
		&friendship.RequestedAt,
		&friendship.StatusChangedAt,
	)

	if err != nil {
		return nil, err
	}

	return friendship, nil
}

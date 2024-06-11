package store

import (
	"encoding/json"
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

func scanFriend(rows Scanner) (*models.Friend, error) {
	friend := &models.Friend{}

	err := rows.Scan(
		&friend.ID,
		&friend.Username,
		&friend.Email,
		&friend.Active,
		&friend.Password,
		&friend.CreatedAt,
		&friend.UpdatedAt,
		&friend.AcceptedAt,
	)

	if err != nil {
		return nil, err
	}

	return friend, nil
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
		&friendship.Seen,
		&friendship.RequestedAt,
		&friendship.StatusChangedAt,
	)

	if err != nil {
		return nil, err
	}

	return friendship, nil
}

// Scans friend request from the given Scanner, this function accepts sql.Rows and sql.Row as Scanner
// and returns a FriendRequest and an error if the scan fails
func scanFriendRequest(s Scanner) (*models.FriendRequest, error) {
	friendRequest := &models.FriendRequest{
		User: &models.User{},
	}

	err := s.Scan(
		&friendRequest.ID,
		&friendRequest.Status,
		&friendRequest.RequestedAt,
		&friendRequest.StatusChangedAt,
		&friendRequest.User.ID,
		&friendRequest.User.Username,
		&friendRequest.User.Email,
		&friendRequest.User.Active,
		&friendRequest.User.Password,
		&friendRequest.User.CreatedAt,
		&friendRequest.User.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return friendRequest, nil
}

func scanChat(scanner Scanner) (*models.Chat, error) {
	chat := &models.Chat{}

	err := scanner.Scan(
		&chat.ID,
		&chat.Name,
		&chat.Type,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func scanMessage(scanner Scanner) (*models.Message, error) {
	message := &models.Message{}
	err := scanner.Scan(
		&message.ID,
		&message.Text,
		&message.CreatedAt,
		&message.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func scanMessageWithUser(scanner Scanner) (*models.MessageWithUser, error) {
	message := &models.MessageWithUser{
		User: &models.User{},
	}
	err := scanner.Scan(
		&message.ID,
		&message.Text,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.User.ID,
		&message.User.Username,
		&message.User.Email,
		&message.User.Active,
		&message.User.Password,
		&message.User.CreatedAt,
		&message.User.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func scanFriendRequestNotification(scanner Scanner) (*models.FriendRequestNotification, error) {
	notificationDto := &models.NotificationDTO{}
	err := scanner.Scan(
		&notificationDto.BaseNotification.Base.ID,
		&notificationDto.BaseNotification.Type,
		&notificationDto.BaseNotification.UserID,
		&notificationDto.Data,
		&notificationDto.BaseNotification.Seen,
		&notificationDto.BaseNotification.Base.CreatedAt,
		&notificationDto.BaseNotification.Base.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	notificationData := &models.FriendRequestNotificationData{}
	err = json.Unmarshal(notificationDto.Data, notificationData)
	if err != nil {
		return nil, err
	}

	return &models.FriendRequestNotification{
		BaseNotification: notificationDto.BaseNotification,
		Data:             notificationData,
	}, nil
}

// Scans only single entry from query that has to be an integer,
// returns id from table or -1 and error when scan returned error
func scanID(scanner Scanner) (int, error) {
	var ID int
	err := scanner.Scan(&ID)
	if err != nil {
		return -1, err
	}
	return ID, nil
}

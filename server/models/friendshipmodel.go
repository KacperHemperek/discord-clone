package models

import (
	"time"
)

type Friendship struct {
	ID              int       `json:"id"`
	InviterID       int       `json:"inviterId"`
	FriendID        int       `json:"friendId"`
	Seen            bool      `json:"seen"`
	Status          string    `json:"status"`
	RequestedAt     time.Time `json:"requestedAt"`
	StatusChangedAt NullTime  `json:"statusChangedAt"`
}

type FriendRequest struct {
	ID              int       `json:"id"`
	User            *User     `json:"user"`
	Status          string    `json:"status"`
	RequestedAt     time.Time `json:"requestedAt"`
	StatusChangedAt NullTime  `json:"statusChangedAt"`
}

type Friend struct {
	AcceptedAt time.Time `json:"acceptedAt"`
	User
}

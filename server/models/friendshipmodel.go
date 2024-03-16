package models

type Friendship struct {
	ID              int    `json:"id"`
	InviterID       int    `json:"inviterId"`
	FriendID        int    `json:"friendId"`
	Accepted        bool   `json:"accepted"`
	Seen            bool   `json:"seen"`
	Status          string `json:"status"`
	RequestedAt     string `json:"requestedAt"`
	StatusChangedAt string `json:"statusChangedAt"`
}

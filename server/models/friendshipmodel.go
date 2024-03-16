package models

type Friendship struct {
	ID              int    `json:"id"`
	UserID          int    `json:"userId"`
	FriendID        int    `json:"friendId"`
	Accepted        bool   `json:"accepted"`
	Seen            bool   `json:"seen"`
	RequestedAt     string `json:"requestedAt"`
	StatusChangedAt string `json:"statusChangedAt"`
}

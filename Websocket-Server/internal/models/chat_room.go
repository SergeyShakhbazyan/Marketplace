package models

import "github.com/gocql/gocql"

type ChatRoom struct {
	ID            gocql.UUID `json:"id,omitempty"`
	CurrentUserID gocql.UUID `json:"currentUserID"`
	TargetUserID  gocql.UUID `json:"targetUserID"`
}

type UserProfile struct {
	UserID      gocql.UUID `json:"userID"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	Avatar      string     `json:"avatar"`
	Status      bool       `json:"status"`
	LastMessage string     `json:"lastMessage"`
	ChatID      gocql.UUID `json:"chatID"`
}

package models

import (
	"github.com/gocql/gocql"
	"github.com/shopspring/decimal"
	"time"
)

type phoneNumber struct {
	CountryCode string `json:"countryCode"`
	Number      string `json:"number"`
}

type User struct {
	UserID       gocql.UUID       `json:"userID"`
	FirstName    string           `json:"firstName"`
	LastName     string           `json:"lastName"`
	Avatar       string           `json:"avatar"`
	PhoneNumber  string           `json:"phoneNumber"`
	Email        string           `json:"email"`
	Password     string           `json:"password"`
	AccountType  string           `json:"accountType"`
	Subscription bool             `json:"subscription"`
	Rating       *decimal.Decimal `json:"rating"`
	CreatedAt    time.Time        `json:"createdAt"`
}

type UserWrapContent struct {
	UserID      gocql.UUID           `json:"userID"`
	FirstName   string               `json:"firstName"`
	LastName    string               `json:"lastName"`
	Avatar      string               `json:"avatar"`
	AccountType string               `json:"accountType"`
	Rating      *decimal.Decimal     `json:"rating"`
	Products    []ProductWrapContent `json:"products,omitempty"`
}

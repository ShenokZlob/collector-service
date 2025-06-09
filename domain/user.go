package domain

import (
	"time"
)

type User struct {
	ID           string              `json:"id"`
	Email        string              `json:"email,omitempty"`
	PasswordHash string              `json:"password_hash,omitempty"`
	TelegramID   int64               `json:"telegram_id,omitempty"`
	FirstName    string              `json:"first_name"`
	LastName     string              `json:"last_name,omitempty"`
	Username     string              `json:"username,omitempty"`
	Collections  []UserCollectionRef `json:"collections,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type UserCollectionRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

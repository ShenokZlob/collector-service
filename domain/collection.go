package domain

import (
	"time"
)

type Collection struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Cards     []Card    `json:"cards,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Card struct {
	ScryfallID string    `json:"scryfall_id"`
	Name       string    `json:"name"`
	CardUrl    string    `json:"card_url"`
	Count      int       `json:"count"`
	AddedAt    time.Time `json:"added_at"`
}

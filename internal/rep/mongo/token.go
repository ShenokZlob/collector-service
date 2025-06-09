package mongorep

import "time"

type TokenInfo struct {
	IDjti     string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	IssuedAt  time.Time `bson:"issued_at"`
	ExpiresAt time.Time `bson:"expires_at"`
	Revoked   bool      `bson:"revoked"`
}

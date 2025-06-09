package mongorep

import (
	"time"

	"github.com/ShenokZlob/collector-service/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	database               = "collector_ouphe_db"
	users_collection       = "users"
	collections_collection = "collections"
	tokens_collection      = "tokens"
)

// user collection
type User struct {
	ObjectID     bson.ObjectID       `bson:"_id,omitempty"`
	Email        string              `bson:"email,omitempty"`
	PasswordHash string              `bson:"password_hash,omitempty"`
	TelegramID   int64               `bson:"telegram_id"`
	FirstName    string              `bson:"first_name"`
	LastName     string              `bson:"last_name,omitempty"`
	Username     string              `bson:"username,omitempty"`
	Collections  []UserCollectionRef `bson:"collections,omitempty"`
	CreatedAt    time.Time           `bson:"created_at,omitempty"`
	UpdatedAt    time.Time           `bson:"updated_at,omitempty"`
}

type UserCollectionRef struct {
	ObjectID bson.ObjectID `bson:"_id,omitempty"`
	Name     string        `bson:"name"`
}

func (u *User) ToDomain() domain.User {
	dCollectionsRef := make([]domain.UserCollectionRef, len(u.Collections))
	for i, v := range u.Collections {
		dCollectionsRef[i] = domain.UserCollectionRef{
			ID:   v.ObjectID.Hex(),
			Name: v.Name,
		}
	}

	return domain.User{
		ID:           u.ObjectID.Hex(),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		TelegramID:   u.TelegramID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Username:     u.Username,
		Collections:  dCollectionsRef,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func UserFromDomain(domainUser domain.User) (User, error) {
	var userObjectID bson.ObjectID
	var err error

	if domainUser.ID != "" {
		userObjectID, err = bson.ObjectIDFromHex(domainUser.ID)
		if err != nil {
			return User{}, err
		}
	}

	collectionsRef := make([]UserCollectionRef, len(domainUser.Collections))
	for i, v := range domainUser.Collections {
		collRefObjectID, err := bson.ObjectIDFromHex(v.ID)
		if err != nil {
			return User{}, err
		}
		collectionsRef[i] = UserCollectionRef{
			ObjectID: collRefObjectID,
			Name:     v.Name,
		}
	}

	return User{
		ObjectID:     userObjectID,
		Email:        domainUser.Email,
		PasswordHash: domainUser.PasswordHash,
		TelegramID:   domainUser.TelegramID,
		FirstName:    domainUser.FirstName,
		LastName:     domainUser.LastName,
		Username:     domainUser.Username,
		Collections:  collectionsRef,
		CreatedAt:    domainUser.CreatedAt,
		UpdatedAt:    domainUser.UpdatedAt,
	}, nil
}

// collection's collection
type Collection struct {
	ObjectID  bson.ObjectID `bson:"_id,omitempty"`
	UserID    bson.ObjectID `bson:"user_id"`
	Name      string        `bson:"name"`
	Cards     []Card        `bson:"cards,omitempty"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

type Card struct {
	ScryfallID string    `bson:"scryfall_id"`
	Name       string    `bson:"name"`
	CardUrl    string    `bson:"card_url"`
	Count      int       `bson:"count"`
	AddedAt    time.Time `bson:"added_at"`
}

func (c *Collection) ToDomain() domain.Collection {
	domainCards := make([]domain.Card, len(c.Cards))
	for i, v := range c.Cards {
		domainCards[i] = domain.Card{
			ScryfallID: v.ScryfallID,
			Name:       v.Name,
			CardUrl:    v.CardUrl,
			Count:      v.Count,
			AddedAt:    v.AddedAt,
		}
	}

	return domain.Collection{
		ID:        c.ObjectID.Hex(),
		UserID:    c.UserID.Hex(),
		Name:      c.Name,
		Cards:     domainCards,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func CollectionFromDomain(domainCollection domain.Collection) (Collection, error) {
	collObjectID, err := bson.ObjectIDFromHex(domainCollection.ID)
	if err != nil {
		return Collection{}, err
	}

	userIdObjectID, err := bson.ObjectIDFromHex(domainCollection.UserID)
	if err != nil {
		return Collection{}, err
	}

	cards := make([]Card, len(domainCollection.Cards))
	for i, v := range domainCollection.Cards {
		cards[i] = Card{
			ScryfallID: v.ScryfallID,
			Name:       v.Name,
			CardUrl:    v.CardUrl,
			Count:      v.Count,
			AddedAt:    v.AddedAt,
		}
	}

	return Collection{
		ObjectID:  collObjectID,
		UserID:    userIdObjectID,
		Name:      domainCollection.Name,
		Cards:     cards,
		CreatedAt: domainCollection.CreatedAt,
		UpdatedAt: domainCollection.UpdatedAt,
	}, nil
}

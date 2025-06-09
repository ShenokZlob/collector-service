package collectorclient

import (
	"context"

	dto "github.com/ShenokZlob/collector-service/pkg/contracts"
)

type CollectorClient interface {
	CollectorClientAuth
	CollectorClientCollections
	CollectorClientCards
}

type CollectorClientAuth interface {
	// TODO
	RegisterUser(reqData *dto.RegisterTelegramRequest) (*dto.RegisterTelegramResponse, error)
}

type CollectorClientCollections interface {
	GetUserCollections(ctx context.Context) ([]dto.Collection, error)
	CreateCollection(ctx context.Context, req *dto.CreateCollectionRequest) (*dto.Collection, error)
	RenameCollection(ctx context.Context, collectionID string, req *dto.RenameCollectionRequest) error
	DeleteCollection(ctx context.Context, collectionID string) error
	GetUsersCollectionByName(ctx context.Context, name string) (*dto.Collection, error)

	// TODO: remove in future
	ListCardsInCollection(ctx context.Context, collectionID string) ([]dto.Card, error)
}

type CollectorClientCards interface {
	ListCardsInCollection(ctx context.Context, collectionID string) ([]dto.Card, error)
	AddCardToCollection(ctx context.Context, collectionID string, card *dto.Card) error
	SetCardCountInCollection(ctx context.Context, collectionID string, card *dto.Card) error
	DeleteCardFromCollection(ctx context.Context, collectionID string, cardIdScryfall string) error
}

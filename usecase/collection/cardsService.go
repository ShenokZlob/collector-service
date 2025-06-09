package collection

import (
	"net/http"

	"github.com/ShenokZlob/collector-service/domain"
	"go.uber.org/zap"
)

type CardsService struct {
	cardsRepository CardsRepositorer
	log             *zap.Logger
}

type CardsRepositorer interface {
	GetCollection(collectionId string) (*domain.Collection, *domain.ResponseErr)
	AddCardToCollection(collectionId string, card *domain.Card) *domain.ResponseErr
	SetCardCountInCollection(collectionId string, card *domain.Card) *domain.ResponseErr
	DeleteCardFromCollection(collectionId string, card *domain.Card) *domain.ResponseErr
}

func NewCardsService(log *zap.Logger, cardsRepository CardsRepositorer) *CardsService {
	return &CardsService{
		cardsRepository: cardsRepository,
		log:             log.With(zap.String("service", "cards")),
	}
}

// ListCardsInCollection retrieves all cards in a collection by its ID.
func (cs CardsService) ListCardsInCollection(collectionId string) ([]domain.Card, *domain.ResponseErr) {
	collection, err := cs.cardsRepository.GetCollection(collectionId)
	if err != nil {
		return nil, err
	}

	if collection == nil {
		return nil, &domain.ResponseErr{
			Status:  http.StatusNotFound,
			Message: "Collection not found",
		}
	}

	// For json serialization, ensure Cards is not nil
	if collection.Cards == nil {
		collection.Cards = []domain.Card{}
	}

	return collection.Cards, nil
}

// AddCardToCollection adds a card to a collection by its ID.
func (cs CardsService) AddCardToCollection(collectionId string, card *domain.Card) *domain.ResponseErr {
	return cs.cardsRepository.AddCardToCollection(collectionId, card)
}

// SetCardCountInCollection updates the count of a card in a collection by its ID.
func (cs CardsService) SetCardCountInCollection(collectionId string, card *domain.Card) *domain.ResponseErr {
	return cs.cardsRepository.SetCardCountInCollection(collectionId, card)
}

// DeleteCardFromCollection removes a card from a collection by its ID.
func (cs CardsService) DeleteCardFromCollection(collectionId string, card *domain.Card) *domain.ResponseErr {
	return cs.cardsRepository.DeleteCardFromCollection(collectionId, card)
}

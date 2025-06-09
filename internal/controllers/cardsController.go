package controllers

import (
	"net/http"

	"github.com/ShenokZlob/collector-service/domain"
	dto "github.com/ShenokZlob/collector-service/pkg/contracts"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CardsController отвечает за работу с коллекциями
// @Tags Cards
// @BasePath /
type CardsController struct {
	log          *zap.Logger
	cardsService CardsServicer
}

type CardsServicer interface {
	ListCardsInCollection(collectionId string) ([]domain.Card, *domain.ResponseErr)
	AddCardToCollection(collectionId string, card *domain.Card) *domain.ResponseErr
	SetCardCountInCollection(collectionId string, card *domain.Card) *domain.ResponseErr
	DeleteCardFromCollection(collectionId string, card *domain.Card) *domain.ResponseErr
}

func NewCardsController(log *zap.Logger, cardsService CardsServicer) *CardsController {
	return &CardsController{
		log:          log.With(zap.String("controller", "cards")),
		cardsService: cardsService,
	}
}

// @Summary     Get user's cards in collection
// @Description Получить список карт из коллекции юзера
// @Tags        Cards
// @Security    BearerAuth
// @Produce     json
// @Success     200 {array} dto.Card
// @Failure     401 {object} dto.ErrorResponse
// @Router      /collections/{id}/cards [get]
func (cc CardsController) ListCardsInCollection(ctx *gin.Context) {
	collectionId := ctx.Param("id")

	cardsList, respErr := cc.cardsService.ListCardsInCollection(collectionId)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	var out []dto.Card
	for _, card := range cardsList {
		out = append(out, dto.Card{
			ScryfallID: card.ScryfallID,
			Name:       card.Name,
			CardUrl:    card.CardUrl,
			Count:      card.Count,
		})
	}

	ctx.JSON(200, out)
}

// @Summary     Add a card to user's collection
// @Description Добавить карту в коллекцию юзера
// @Tags        Cards
// @Security    BearerAuth
// @Produce     json
// @Success     201 "No Content"
// @Failure     401 {object} dto.ErrorResponse
// @Router      /collections/{id}/cards [post]
func (cc CardsController) AddCardToCollection(ctx *gin.Context) {
	collectionId := ctx.Param("id")
	var card domain.Card
	if err := ctx.ShouldBindJSON(&card); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	respErr := cc.cardsService.AddCardToCollection(collectionId, &card)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ctx.Status(http.StatusCreated)
}

// @Summary     Set card count in user's collection
// @Description Установить количество карт в коллекции юзера
// @Tags        Cards
// @Security    BearerAuth
// @Produce     json
// @Success     204 "No Content"
// @Failure     401 {object} dto.ErrorResponse
// @Router      /collections/{id}/cards/{card_id} [patch]
func (cc CardsController) SetCardCountInCollection(ctx *gin.Context) {
	collectionId := ctx.Param("id")
	scryfallId := ctx.Param("card_id")

	var card domain.Card
	if err := ctx.ShouldBindJSON(&card); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	card.ScryfallID = scryfallId
	respErr := cc.cardsService.SetCardCountInCollection(collectionId, &card)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary     Delete the card from user's collection
// @Description Удалить карту из коллекции юзера
// @Tags        Cards
// @Security    BearerAuth
// @Produce     json
// @Success     204 "No Content"
// @Failure     401 {object} dto.ErrorResponse
// @Router      /collections/{id}/cards/{card_id} [delete]
func (cc CardsController) DeleteCardFromCollection(ctx *gin.Context) {
	collectionId := ctx.Param("id")
	scryfallId := ctx.Param("card_id")

	respErr := cc.cardsService.DeleteCardFromCollection(collectionId, &domain.Card{ScryfallID: scryfallId})
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

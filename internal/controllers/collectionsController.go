// internal/app/controllers/collections_controller.go
package controllers

import (
	"net/http"

	"github.com/ShenokZlob/collector-service/domain"
	dto "github.com/ShenokZlob/collector-service/pkg/contracts"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CollectionsController отвечает за работу с коллекциями
// @Tags Collections
// @BasePath /
type CollectionsController struct {
	collectionsService CollectionsServicer
	log                *zap.Logger
}

type CollectionsServicer interface {
	GetAll(userId string) ([]domain.UserCollectionRef, *domain.ResponseErr)
	Get(collectionID string) (*domain.Collection, *domain.ResponseErr)
	Create(collection *domain.Collection) (*domain.Collection, *domain.ResponseErr)
	Rename(collection *domain.Collection) (*domain.Collection, *domain.ResponseErr)
	Delete(userID, collectionID string) *domain.ResponseErr
}

// NewCollectionsController создает контроллер коллекций
func NewCollectionsController(log *zap.Logger, collectionsService CollectionsServicer) *CollectionsController {
	return &CollectionsController{
		log:                log.With(zap.String("controller", "collections")),
		collectionsService: collectionsService,
	}
}

// @Summary     Get user's collections
// @Description Получить список коллекций текущего пользователя
// @Tags        Collections
// @Security    BearerAuth
// @Produce     json
// @Success     200 {array} dto.Collection
// @Failure     401 {object} dto.ErrorResponse
// @Router      /collections [get]
func (cc CollectionsController) GetAll(ctx *gin.Context) {
	userID, respErr := getUserFromCtx(ctx)
	if respErr != nil {
		cc.log.Error("GetAllCollections: failed to get userID", zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	cc.log.Info("GetAllCollections: started", zap.String("userID", userID))

	list, respErr := cc.collectionsService.GetAll(userID)
	if respErr != nil {
		cc.log.Error("GetAllCollections: failed to get user's collections",
			zap.String("userID", userID), zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	var out []dto.Collection
	for _, c := range list {
		out = append(out, dto.Collection{ID: c.ID, Name: c.Name})
	}
	cc.log.Info("GetAllCollections: success", zap.String("userID", userID), zap.Int("collections_count", len(out)))
	ctx.JSON(http.StatusOK, out)
}

// @Summary     Get one user's collection by ID
// @Description Получить коллекцию пользователя
// @Tags        Collections
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} dto.Collection
// @Failure     401 {object} dto.ErrorResponse
// @Router      /collections/{id} [get]
func (cc CollectionsController) Get(ctx *gin.Context) {
	userID, respErr := getUserFromCtx(ctx)
	if respErr != nil {
		cc.log.Error("GetCollection: failed to get userID", zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	cc.log.Info("GetCollection: started", zap.String("userID", userID))

	collectionID := ctx.Param("id")
	collection, respErr := cc.collectionsService.Get(collectionID)
	if respErr != nil {
		cc.log.Error("GetlCollection: failed to get collection", zap.String("userID", userID),
			zap.String("collectionID", collectionID), zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	cc.log.Info("GetCollection: success", zap.String("userID", userID), zap.String("collectionID", collectionID))
	ctx.JSON(http.StatusOK, collection)
}

// @Summary     Create new collection
// @Description Создать новую коллекцию с указанным именем
// @Tags        Collections
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       input body dto.CreateCollectionRequest true "Название новой коллекции"
// @Success     201 {object} dto.Collection
// @Failure     400,401 {object} dto.ErrorResponse
// @Router      /collections [post]
func (cc CollectionsController) Create(ctx *gin.Context) {
	userID, respErr := getUserFromCtx(ctx)
	if respErr != nil {
		cc.log.Error("CreateCollection: failed to get userID", zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	cc.log.Info("CreateCollection: started", zap.String("userID", userID))

	var req dto.CreateCollectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		cc.log.Error("CreateCollection: failed to get request body", zap.String("userID", userID), zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	collection := &domain.Collection{UserID: userID, Name: req.Name}
	created, respErr := cc.collectionsService.Create(collection)
	if respErr != nil {
		cc.log.Error("CreateCollection: failed to create collection", zap.String("userID", userID), zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	out := dto.Collection{ID: created.ID, Name: created.Name}
	cc.log.Info("CreateCollection: success", zap.String("userID", userID))
	ctx.JSON(http.StatusCreated, out)
}

// @Summary     Rename collection
// @Description Переименовать коллекцию по ID
// @Tags        Collections
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id   path string                         true "Collection ID"
// @Param       input body dto.RenameCollectionRequest true "Новое имя коллекции"
// @Success     204 {object} dto.Collection
// @Failure     400,401,404 {object} dto.ErrorResponse
// @Router      /collections/{id} [patch]
func (cc CollectionsController) Rename(ctx *gin.Context) {
	userID, respErr := getUserFromCtx(ctx)
	if respErr != nil {
		cc.log.Error("RenameCollection: failed to get userID", zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	cc.log.Info("RenameCollection: started", zap.String("userID", userID))

	collectionID := ctx.Param("id")
	var req dto.RenameCollectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		cc.log.Error("RenameCollection: failed to get request body", zap.String("userID", userID), zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	collection := &domain.Collection{ID: collectionID, UserID: userID, Name: req.Name}
	updatedCollection, respErr := cc.collectionsService.Rename(collection)
	if respErr != nil {
		cc.log.Error("RenameCollection: failed to rename collection", zap.String("userID", userID), zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	out := dto.Collection{ID: updatedCollection.ID, Name: updatedCollection.Name}
	cc.log.Info("RenameCollection: success", zap.String("userID", userID))
	ctx.JSON(http.StatusNoContent, out)
}

// @Summary     Delete collection
// @Description Удалить коллекцию по ID
// @Tags        Collections
// @Security    BearerAuth
// @Produce     json
// @Param       id path string true "Collection ID"
// @Success     204 "No Content"
// @Failure     401,404 {object} dto.ErrorResponse
// @Router      /collections/{id} [delete]
func (cc CollectionsController) Delete(ctx *gin.Context) {
	userID, respErr := getUserFromCtx(ctx)
	if respErr != nil {
		cc.log.Error("Deletecollection: failed to get userID", zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	cc.log.Info("DeleteCollection: started", zap.String("userID", userID))

	collectionID := ctx.Param("id")
	respErr = cc.collectionsService.Delete(userID, collectionID)
	if respErr != nil {
		cc.log.Error("DeleteCollection: failed to delete collection", zap.String("userID", userID), zap.String("collectionID", collectionID), zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	cc.log.Info("DeleteCollection: success", zap.String("userID", userID), zap.String("collectionID", collectionID))
	ctx.Status(http.StatusNoContent)
}

func getUserFromCtx(ctx *gin.Context) (string, *domain.ResponseErr) {
	val, ok := ctx.Get("userID")
	if !ok {
		return "", &domain.ResponseErr{Status: http.StatusUnauthorized, Message: "Don't have user ID"}
	}
	userID, ok := val.(string)
	if !ok || userID == "" {
		return "", &domain.ResponseErr{Status: http.StatusUnauthorized, Message: "Invalid user ID type"}
	}
	return userID, nil
}

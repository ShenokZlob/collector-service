package collection

import (
	"net/http"
	"regexp"

	"github.com/ShenokZlob/collector-service/domain"
	"go.uber.org/zap"
)

type CollectionsService struct {
	log                  *zap.Logger
	collectionRepository CollectionsRepositorer
}

type CollectionsRepositorer interface {
	// GetAllUsersCollections(userId string) ([]*domain.UserCollectionRef, *domain.ResponseErr)
	GetUser(userId string) (*domain.User, *domain.ResponseErr)
	GetCollection(collectionID string) (*domain.Collection, *domain.ResponseErr)
	CreateCollection(collection *domain.Collection) (*domain.Collection, *domain.ResponseErr)
	RenameCollection(collection *domain.Collection) (*domain.Collection, *domain.ResponseErr)
	DeleteCollection(userID, collectionID string) *domain.ResponseErr
}

func NewCollectionsService(log *zap.Logger, collectionRepository CollectionsRepositorer) *CollectionsService {
	return &CollectionsService{
		log:                  log.With(zap.String("usecase", "collections")),
		collectionRepository: collectionRepository,
	}
}

func (cs CollectionsService) GetAll(userID string) ([]domain.UserCollectionRef, *domain.ResponseErr) {
	user, respErr := cs.collectionRepository.GetUser(userID)
	if respErr != nil {
		cs.log.Error("Failed to find user", zap.String("userID", userID))
		return nil, respErr
	}
	return user.Collections, nil
}

func (cs CollectionsService) Get(collectionID string) (*domain.Collection, *domain.ResponseErr) {
	if !isValidCollectionID(collectionID) {
		cs.log.Warn("Invalid collection ID", zap.String("collectionID", collectionID))
		return nil, &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID",
		}
	}

	return cs.collectionRepository.GetCollection(collectionID)
}

func (cs CollectionsService) Create(collection *domain.Collection) (*domain.Collection, *domain.ResponseErr) {
	if !isValidCollectionID(collection.ID) {
		cs.log.Warn("Invalid collection ID", zap.String("collectionID", collection.ID))
		return nil, &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID",
		}
	}

	if !isValidCollectionName(collection.Name) {
		cs.log.Warn("Invalid collection name", zap.String("collectionName", collection.Name))
		return nil, &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection name",
		}
	}

	return cs.collectionRepository.CreateCollection(collection)
}

func (cs CollectionsService) Rename(collection *domain.Collection) (*domain.Collection, *domain.ResponseErr) {
	if !isValidCollectionID(collection.ID) {
		cs.log.Warn("Invalid collection ID", zap.String("collectionID", collection.ID))
		return nil, &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID",
		}
	}

	if !isValidCollectionName(collection.Name) {
		cs.log.Warn("Invalid collection name", zap.String("collectionName", collection.Name))
		return nil, &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection name",
		}
	}

	return cs.collectionRepository.RenameCollection(collection)
}

func (cs CollectionsService) Delete(userID, collectionID string) *domain.ResponseErr {
	if !isValidCollectionID(collectionID) {
		cs.log.Warn("Invalid collection ID", zap.String("collectionID", collectionID))
		return &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID",
		}
	}

	return cs.collectionRepository.DeleteCollection(userID, collectionID)
}

var oidRegexp = regexp.MustCompile("^[0-9a-fA-F]{24}$")

func isValidCollectionID(collecionID string) bool {
	return oidRegexp.MatchString(collecionID)
}

var nameRegexp = regexp.MustCompile("^.{1,20}$")

func isValidCollectionName(name string) bool {
	return nameRegexp.MatchString(name)
}

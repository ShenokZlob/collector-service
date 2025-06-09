package mongorep

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ShenokZlob/collector-service/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	client *mongo.Client
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		client: client,
	}
}

// CreateUser
func (r Repository) CreateUser(domainUser *domain.User) (*domain.User, *domain.ResponseErr) {
	user, err := UserFromDomain(*domainUser)
	if err != nil {
		return nil, &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	storage := r.client.Database(database).Collection(users_collection)
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	result, err := storage.InsertOne(context.TODO(), user)
	if err != nil {
		var we mongo.WriteException
		if errors.As(err, &we) {
			for _, e := range we.WriteErrors {
				if e.Code == 11000 { // Duplicate key error
					return nil, &domain.ResponseErr{
						Status:  http.StatusConflict,
						Message: "User with this Telegram ID already exists",
					}
				}
			}
		}
		return nil, &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	insertedID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil, &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get inserted user ID",
		}
	}

	var createdUser User
	filter := bson.M{"_id": insertedID}
	err = storage.FindOne(context.TODO(), filter).Decode(&createdUser)
	if err != nil {
		return nil, &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	domainCreatedUser := createdUser.ToDomain()
	return &domainCreatedUser, nil
}

// GetUser
func (r Repository) GetUser(userId string) (*domain.User, *domain.ResponseErr) {
	objectID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID format",
		}
	}

	collection := r.client.Database(database).Collection(users_collection)
	filter := bson.D{{Key: "_id", Value: objectID}}

	var user User
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &domain.ResponseErr{
				Status:  http.StatusNotFound,
				Message: "User not found",
			}
		}
		return nil, &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Find user error: %v", err),
		}
	}

	domainUser := user.ToDomain()
	return &domainUser, nil
}

// FindByEmail find user in user's collection by email
func (r Repository) FindByEmail(email string) (*domain.User, *domain.ResponseErr) {
	collection := r.client.Database(database).Collection(users_collection)
	filter := bson.M{"email": email}

	var user User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &domain.ResponseErr{
				Status:  http.StatusNotFound,
				Message: "User not found",
			}
		}
		return nil, &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Find user error: %v", err),
		}
	}

	domainUser := user.ToDomain()
	return &domainUser, nil
}

func (r Repository) AddToken(userID, jti string, issued, expires time.Time) *domain.ResponseErr {
	storage := r.client.Database(database).Collection(tokens_collection)
	token := TokenInfo{
		IDjti:     jti,
		UserID:    userID,
		IssuedAt:  issued,
		ExpiresAt: expires,
		Revoked:   false,
	}

	_, err := storage.InsertOne(context.TODO(), &token)
	if err != nil {
		var we mongo.WriteException
		if errors.As(err, &we) {
			for _, e := range we.WriteErrors {
				if e.Code == 11000 { // Duplicate key error
					return &domain.ResponseErr{
						Status:  http.StatusConflict,
						Message: "User with this Telegram ID already exists",
					}
				}
			}
		}
		return &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

func (r Repository) AddToBlackList(jti string) *domain.ResponseErr {
	storage := r.client.Database(database).Collection(tokens_collection)
	filter := bson.M{"_id": jti}
	update := bson.M{
		"$set": bson.M{"revoked": true},
	}

	_, err := storage.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &domain.ResponseErr{
				Status:  http.StatusNotFound,
				Message: "User not found",
			}
		}
		return &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Find user error: %v", err),
		}
	}

	return nil
}

// FindUserByTelegramID searches for a user by their Telegram ID
func (r Repository) FindUserByTelegramID(telegramId int64) (*domain.User, *domain.ResponseErr) {
	collection := r.client.Database(database).Collection(users_collection)
	filter := bson.D{{Key: "telegram_id", Value: telegramId}}

	var user User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &domain.ResponseErr{
				Status:  http.StatusNotFound,
				Message: "User not found",
			}
		}
		fmt.Println(user)
		return nil, &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Find user error: %v", err),
		}
	}

	domainUser := user.ToDomain()
	return &domainUser, nil
}

func (r Repository) CreateCollection(domainCollection *domain.Collection) (*domain.Collection, *domain.ResponseErr) {
	collection, err := CollectionFromDomain(*domainCollection)
	if err != nil {
		return nil, &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	ctx := context.TODO()
	session, err := r.client.StartSession()
	if err != nil {
		return nil, &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Failed to start session: %v", err),
		}
	}
	defer session.EndSession(ctx)

	// Start session (transaction)
	var domainCreatedCollection domain.Collection
	err = mongo.WithSession(ctx, session, func(ctx context.Context) error {
		// Add collection to collections_collection
		storage := r.client.Database(database).Collection(collections_collection)
		collection.CreatedAt = time.Now()
		collection.UpdatedAt = collection.CreatedAt

		result, err := storage.InsertOne(ctx, collection)
		if err != nil {
			return &domain.ResponseErr{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		// Get created collection
		insertedID, ok := result.InsertedID.(bson.ObjectID)
		if !ok {
			return &domain.ResponseErr{
				Status:  http.StatusInternalServerError,
				Message: "Created collection has invalid ID",
			}
		}
		var createdCollection Collection
		filter := bson.M{"_id": insertedID}
		err = storage.FindOne(ctx, filter).Decode(&createdCollection)
		if err != nil {
			return &domain.ResponseErr{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		// Add created collection to collections_users
		storage = r.client.Database(database).Collection(users_collection)
		filter = bson.M{"_id": createdCollection.UserID}
		update := bson.D{
			{Key: "$push", Value: bson.D{{Key: "collections", Value: UserCollectionRef{
				ObjectID: collection.ObjectID,
				Name:     collection.Name,
			}}}},
			{Key: "$set", Value: bson.D{{Key: "updated_at", Value: time.Now()}}},
		}
		_, err = storage.UpdateOne(ctx, filter, update)
		if err != nil {
			return &domain.ResponseErr{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		domainCreatedCollection = createdCollection.ToDomain()
		return nil
	})
	if err != nil {
		if e, ok := err.(*domain.ResponseErr); ok {
			return nil, e
		}
		return nil, &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return &domainCreatedCollection, nil
}

func (r Repository) RenameCollection(domainCollection *domain.Collection) (*domain.Collection, *domain.ResponseErr) {
	collection, err := CollectionFromDomain(*domainCollection)
	if err != nil {
		return nil, &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID format",
		}
	}

	ctx := context.TODO()
	session, err := r.client.StartSession()
	if err != nil {
		return nil, &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Failed to start session: %v", err),
		}
	}
	defer session.EndSession(ctx)

	// Start session (transaction)
	var domainRenamedCollection domain.Collection
	err = mongo.WithSession(ctx, session, func(ctx context.Context) error {
		// Rename collection in collections_collection
		storage := r.client.Database(database).Collection(collections_collection)
		filter := bson.M{"_id": collection.ObjectID}
		update := bson.M{
			"$set": bson.M{
				"name":       collection.Name,
				"updated_at": time.Now(),
			},
		}
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		var updatedCollection Collection
		err = storage.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedCollection)
		if err != nil {
			return &domain.ResponseErr{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("Failed to find collection: %v", err),
			}
		}

		// Update collection name in user's collections
		storage = r.client.Database(database).Collection(users_collection)
		filter = bson.M{"collections._id": updatedCollection.ObjectID}
		update = bson.M{
			"$set": bson.M{
				"collections.$.name": collection.Name,
				"updated_at":         time.Now(),
			},
		}
		_, err = storage.UpdateOne(ctx, filter, update)
		if err != nil {
			return &domain.ResponseErr{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		domainRenamedCollection = updatedCollection.ToDomain()
		return nil
	})
	if err != nil {
		if e, ok := err.(*domain.ResponseErr); ok {
			return nil, e
		}
		return nil, &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return &domainRenamedCollection, nil
}

func (r Repository) DeleteCollection(userID, collectionID string) *domain.ResponseErr {
	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID format",
		}
	}

	collectionObjectID, err := bson.ObjectIDFromHex(collectionID)
	if err != nil {
		return &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID format",
		}
	}

	ctx := context.TODO()
	session, err := r.client.StartSession()
	if err != nil {
		return &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Failed to start session: %v", err),
		}
	}
	defer session.EndSession(ctx)

	// Start session (transaction)
	err = mongo.WithSession(ctx, session, func(ctx context.Context) error {
		// Delete collection from collections_collection
		storage := r.client.Database(database).Collection(collections_collection)
		filter := bson.M{"_id": collectionObjectID}
		_, err = storage.DeleteOne(ctx, filter)
		if err != nil {
			return &domain.ResponseErr{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("Delete collection error: %v", err),
			}
		}

		// Delete collection from user's collections
		storage = r.client.Database(database).Collection(users_collection)
		filter = bson.M{"_id": userObjectID}
		update := bson.M{
			"$pull": bson.M{
				"collections": bson.M{"_id": collectionObjectID},
			},
			"$set": bson.M{"updated_at": time.Now()},
		}

		_, err = storage.UpdateOne(ctx, filter, update)
		if err != nil {
			return &domain.ResponseErr{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("Error updating user collections: %v", err),
			}
		}

		return nil
	})
	if err != nil {
		if e, ok := err.(*domain.ResponseErr); ok {
			return e
		}
		return &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// GetCollection gets all information about collection by ID
func (r Repository) GetCollection(collectionId string) (*domain.Collection, *domain.ResponseErr) {
	collObjectID, err := bson.ObjectIDFromHex(collectionId)
	if err != nil {
		return nil, &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID format",
		}
	}
	storage := r.client.Database(database).Collection(collections_collection)
	filter := bson.M{"_id": collObjectID}

	var collection Collection
	err = storage.FindOne(context.TODO(), filter).Decode(&collection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &domain.ResponseErr{
				Status:  http.StatusNotFound,
				Message: "Collection not found",
			}
		}
		return nil, &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Find collection error: %v", err),
		}
	}

	domainCollection := collection.ToDomain()
	return &domainCollection, nil
}

func (r Repository) AddCardToCollection(collectionId string, card *domain.Card) *domain.ResponseErr {
	// TODO: need to use transaction here
	// TODO: replace collectionId on ObjectID in service layer
	objectId, err := bson.ObjectIDFromHex(collectionId)
	if err != nil {
		return &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID format",
		}
	}

	collection := r.client.Database(database).Collection(collections_collection)
	card.AddedAt = time.Now()

	// Try to update the card count first
	filter := bson.M{
		"_id":               objectId,
		"cards.scryfall_id": card.ScryfallID,
	}
	update := bson.M{
		"$inc": bson.M{"cards.$.count": card.Count},
		"$set": bson.M{"updated_at": time.Now()},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Error updating card count: %v", err),
		}
	}

	if result.MatchedCount > 0 {
		return nil
	}

	// If the card doesn't exist, add it to the collection
	filter = bson.M{"_id": objectId}
	update = bson.M{
		"$push": bson.M{"cards": card},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Error adding new card: %v", err),
		}
	}

	return nil
}

func (r Repository) SetCardCountInCollection(collectionId string, card *domain.Card) *domain.ResponseErr {
	// TODO: need to use transaction here
	// TODO: replace collectionId on ObjectID in service layer
	objectId, err := bson.ObjectIDFromHex(collectionId)
	if err != nil {
		return &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID format",
		}
	}

	collection := r.client.Database(database).Collection(collections_collection)
	filter := bson.D{
		{Key: "_id", Value: objectId},
		{Key: "cards.scryfall_id", Value: card.ScryfallID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "cards.$.count", Value: card.Count},
			{Key: "updated_at", Value: time.Now()},
		}},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	// var updatedColection Collection
	res := collection.FindOneAndUpdate(context.TODO(), filter, update, opts)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return &domain.ResponseErr{
				Status:  http.StatusNotFound,
				Message: "Collection not found",
			}
		}
		return &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Update collection error: %v", err),
		}
	}

	return nil
}

func (r Repository) DeleteCardFromCollection(collectionId string, card *domain.Card) *domain.ResponseErr {
	// TODO: replace collectionId on ObjectID in service layer
	objectId, err := bson.ObjectIDFromHex(collectionId)
	if err != nil {
		return &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID format",
		}
	}

	collection := r.client.Database(database).Collection(collections_collection)
	filter := bson.D{
		{Key: "_id", Value: objectId},
		{Key: "cards.scryfall_id", Value: card.ScryfallID},
	}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "cards", Value: bson.D{
				{Key: "scryfall_id", Value: card.ScryfallID},
			}},
		}},
		{Key: "$set", Value: bson.D{
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	res := collection.FindOneAndUpdate(context.TODO(), filter, update)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return &domain.ResponseErr{
				Status:  http.StatusNotFound,
				Message: "Collection not found",
			}
		}
		return &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Update collection error: %v", err),
		}
	}

	return nil
}

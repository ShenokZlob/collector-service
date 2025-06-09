package collectorclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShenokZlob/collector-service/pkg/authctx"
	dto "github.com/ShenokZlob/collector-service/pkg/contracts"
	"go.uber.org/zap"
)

type HTTPCollectorClient struct {
	URL        string
	Log        *zap.Logger
	ClientHTTP *http.Client
}

func NewHTTPCollectorClient(url string, log *zap.Logger) *HTTPCollectorClient {
	return &HTTPCollectorClient{
		URL:        url,
		Log:        log,
		ClientHTTP: http.DefaultClient,
	}
}

// RegisterUser reg the user in collection service
func (c *HTTPCollectorClient) RegisterUser(reqdata *dto.RegisterTelegramRequest) (*dto.RegisterTelegramResponse, error) {
	c.Log.Info("Registering user in collector service", zap.String("method", "HTTPCollectorClient.RegisterUser"), zap.Int("telegram_id", int(reqdata.TelegramID)))

	body, err := json.Marshal(reqdata)
	if err != nil {
		c.Log.Error("Failed to marshal request data", zap.Error(err))
		return nil, err
	}

	resp, err := http.Post(c.URL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.Log.Error("Failed to send request to collector service", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		c.Log.Error("Failed to register user in collector service", zap.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("failed to register user in collector service, status code: %d", resp.StatusCode)
	}

	var respData dto.RegisterTelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		c.Log.Error("Failed to decode response data", zap.Error(err))
		return nil, err
	}

	return &respData, nil
}

// GetCollections gets list of collections for user
// Need JWT token for this opperation
// Authorization: Bearer TOKEN
func (c *HTTPCollectorClient) GetUserCollections(ctx context.Context) ([]dto.Collection, error) {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return nil, fmt.Errorf("authorization token is missing")
	}

	c.Log.Info("Get user's list of collections", zap.String("method", "HTTPCollectorClient.GetUserCollections"), zap.String("token_auth", token))

	request, err := http.NewRequest(http.MethodGet, c.URL+"/collections", nil)
	if err != nil {
		c.Log.Error("Failed to create request", zap.Error(err))
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)
	c.Log.Debug("Req token", zap.String("header", request.Header.Get("Authorization")))

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorRepsonse dto.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorRepsonse); err != nil {
			c.Log.Error("Failed to decode error response", zap.Error(err))
			return nil, fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to get user's collections", zap.String("message", errorRepsonse.Message))
		return nil, fmt.Errorf("failed to get user's collections, status code: %d", resp.StatusCode)
	}

	var collections []dto.Collection
	err = json.NewDecoder(resp.Body).Decode(&collections)
	if err != nil {
		c.Log.Error("Failed to decode a body request", zap.Error(err))
		return nil, err
	}

	return collections, nil
}

// Need JWT token for this opperation
func (c *HTTPCollectorClient) CreateCollection(ctx context.Context, req *dto.CreateCollectionRequest) (*dto.Collection, error) {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return nil, fmt.Errorf("authorization token is missing")
	}

	// doesn't look good to logging this information!!!
	c.Log.Info("Create collection", zap.String("metod", "HTTPCollectorClient.CreateCollection"), zap.String("token_auth", token))

	body, err := json.Marshal(&req)
	if err != nil {
		c.Log.Error("Failed to marshal data", zap.Error(err))
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, c.URL+"/collections", bytes.NewBuffer(body))
	if err != nil {
		c.Log.Error("Failed to create request", zap.Error(err))
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to send request to collector service", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResponse dto.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			c.Log.Error("Failed to decode error response", zap.Error(err))
			return nil, fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to create collection's for user", zap.String("message", errorResponse.Message))
		return nil, fmt.Errorf("failed to create collection's, status code: %d", resp.StatusCode)
	}

	var collection dto.Collection
	err = json.NewDecoder(resp.Body).Decode(&collection)
	if err != nil {
		c.Log.Error("Failed to decode a responser body", zap.Error(err))
		return nil, err
	}

	return &collection, nil
}

// Need JWT token for this opperation
func (c *HTTPCollectorClient) RenameCollection(ctx context.Context, collectionID string, req *dto.RenameCollectionRequest) error {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return fmt.Errorf("authorization token is missing")
	}

	c.Log.Info("Rename collection", zap.String("collection_id", collectionID), zap.String("method", "HTTPCollectorClient.RenameCollection"), zap.String("token_auth", token))

	body, err := json.Marshal(req)
	if err != nil {
		c.Log.Error("Failed to marshal request data")
		return err
	}

	request, err := http.NewRequest(http.MethodPatch, c.URL+"/collections/"+collectionID, bytes.NewBuffer(body))
	if err != nil {
		c.Log.Error("Failed to create request", zap.Error(err))
		return err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		var errorResponse dto.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			c.Log.Error("Failed to decode error response", zap.Error(err))
			return fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to rename collection", zap.String("message", errorResponse.Message))
		return fmt.Errorf("failed to rename the user's collection, status code: %d", resp.StatusCode)
	}

	return nil
}

// Need JWT token for this opperation
func (c *HTTPCollectorClient) DeleteCollection(ctx context.Context, collectionID string) error {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return fmt.Errorf("authorization token is missing")
	}

	// doesn't look good to logging this information!!!
	c.Log.Info("Delete collection", zap.String("method", "HTTPCollectorClient.DeleteCollection"), zap.String("token_auth", token))

	request, err := http.NewRequest(http.MethodDelete, c.URL+"/collections/"+collectionID, nil)
	if err != nil {
		c.Log.Error("Failed to prepare a request", zap.Error(err))
		return err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		var errorResponse dto.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			c.Log.Error("Failed to decode error response", zap.Error(err))
			return fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to delete the user's collection", zap.String("message", errorResponse.Message))
		return fmt.Errorf("failed to delete the user's collection, status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *HTTPCollectorClient) GetUsersCollectionByName(ctx context.Context, collectionName string) (*dto.Collection, error) {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return nil, fmt.Errorf("authorization token is missing")
	}

	c.Log.Info("Get user's collection by name", zap.String("method", "HTTPCollectorClient.GetUsersCollectionByName"), zap.String("token_auth", token), zap.String("collection_name", collectionName))

	request, err := http.NewRequest(http.MethodGet, c.URL+"/collections/name/"+collectionName, nil)
	if err != nil {
		c.Log.Error("Failed to create request", zap.Error(err))
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse dto.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			c.Log.Error("Failed to decode error response", zap.Error(err))
			return nil, fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to get user's collection by name", zap.String("message", errorResponse.Message))
		return nil, fmt.Errorf("failed to get user's collection by name, status code: %d", resp.StatusCode)
	}

	var collection dto.Collection
	err = json.NewDecoder(resp.Body).Decode(&collection)
	if err != nil {
		c.Log.Error("Failed to decode a body request", zap.Error(err))
		return nil, err
	}

	return &collection, nil
}

func (c *HTTPCollectorClient) ListCardsInCollection(ctx context.Context, collectionID string) ([]dto.Card, error) {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return nil, fmt.Errorf("authorization token is missing")
	}

	c.Log.Info("Get cards from collection", zap.String("method", "HTTPCollectorClient.ListCardsInCollection"), zap.String("token_auth", token), zap.String("collection_id", collectionID))

	request, err := http.NewRequest(http.MethodGet, c.URL+fmt.Sprintf("/collections/%s/cards", collectionID), nil)
	if err != nil {
		c.Log.Error("Failed to create request", zap.Error(err))
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse dto.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			c.Log.Error("Failed to decode error response", zap.Error(err))
			return nil, fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to get cards from collection", zap.String("message", errorResponse.Message))
		return nil, fmt.Errorf("failed to get cards from collection, status code: %d", resp.StatusCode)
	}

	var cardsList []dto.Card
	err = json.NewDecoder(resp.Body).Decode(&cardsList)
	if err != nil {
		c.Log.Error("Failed to decode a body request", zap.Error(err))
		return nil, err
	}

	return cardsList, nil
}

func (c *HTTPCollectorClient) AddCardToCollection(ctx context.Context, collectionID string, card *dto.Card) error {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return fmt.Errorf("authorization token is missing")
	}

	c.Log.Info("Add card to collection", zap.String("method", "HTTPCollectorClient.AddCardToCollection"), zap.String("token_auth", token), zap.String("collection_id", collectionID))

	body, err := json.Marshal(card)
	if err != nil {
		c.Log.Error("Failed to marshal card data", zap.Error(err))
		return err
	}

	request, err := http.NewRequest(http.MethodPost, c.URL+fmt.Sprintf("/collections/%s/cards", collectionID), bytes.NewBuffer(body))
	if err != nil {
		c.Log.Error("Failed to create request", zap.Error(err))
		return err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResponse dto.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			c.Log.Error("Failed to decode error response", zap.Error(err))
			return fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to add card to collection", zap.String("message", errorResponse.Message))
		return fmt.Errorf("failed to add card to collection, status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *HTTPCollectorClient) SetCardCountInCollection(ctx context.Context, collectionID string, card *dto.Card) error {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return fmt.Errorf("authorization token is missing")
	}

	c.Log.Info("Set card's count in collection", zap.String("method", "HTTPCollectorClient.SetCardCountInCollection"), zap.String("token_auth", token), zap.String("collection_id", collectionID))

	request, err := http.NewRequest(http.MethodPatch, c.URL+fmt.Sprintf("/collections/%s/cards/%s", collectionID, card.ScryfallID), nil)
	if err != nil {
		c.Log.Error("Failed to create request", zap.Error(err))
		return err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		var errorResponse dto.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			c.Log.Error("Failed to decode error response", zap.Error(err))
			return fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to set card's count in collection", zap.String("message", errorResponse.Message))
		return fmt.Errorf("failed to set card's count in collection, status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *HTTPCollectorClient) DeleteCardFromCollection(ctx context.Context, collectionID string, cardIdScryfall string) error {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return fmt.Errorf("authorization token is missing")
	}

	c.Log.Info("Delete card from collection", zap.String("method", "HTTPCollectorClient.DeleteCardFromCollection"),
		zap.String("token_auth", token), zap.String("collection_id", collectionID), zap.String("card_id_scryfall", cardIdScryfall))

	request, err := http.NewRequest(http.MethodGet, c.URL+fmt.Sprintf("/collections/%s/cards/%s", collectionID, cardIdScryfall), nil)
	if err != nil {
		c.Log.Error("Failed to create request", zap.Error(err))
		return err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		var errorResponse dto.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			c.Log.Error("Failed to decode error response", zap.Error(err))
			return fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to delete card from collection", zap.String("message", errorResponse.Message))
		return fmt.Errorf("failed to delete card from collection, status code: %d", resp.StatusCode)
	}

	return nil
}

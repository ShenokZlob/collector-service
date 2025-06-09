package dto

// CreateCardRequest — запрос для добавления новой карты в коллекцию
// @Description Запрос для добавления карты с Scryfall ID, именем и URL изображения
// @example { "scryfall_id": "12345678-1234-1234-1234-123456789012", "name": "Black Lotus", "card_url": "https://example.com/black-lotus.jpg", "count": 1 }
type AddCardRequest struct {
	ScryfallID string `json:"scryfall_id" binding:"required" example:"12345678-1234-1234-1234-123456789012"`
	Name       string `json:"name" binding:"required" example:"Black Lotus"`
	CardUrl    string `json:"card_url" binding:"required" example:"https://example.com/black-lotus.jpg"`
	Count      int    `json:"count" binding:"required" example:"1"`
}

// SetCardCountRequest — запрос для изменения количества карт в коллекции
// @Description Запрос для изменения количества карт в коллекции по Scryfall ID
// @example { "scryfall_id": "12345678-1234-1234-1234-123456789012", "count": 2 }
type SetCardsCountRequest struct {
	ScryfallID string `json:"scryfall_id" binding:"required" example:"12345678-1234-1234-1234-123456789012"`
	Count      int    `json:"count" binding:"required" example:"2"` // Новое количество карт
}

// Card - модель карты в ответах
// @Description Модель карты с Scryfall ID, именем, URL изображения и количеством
// @example { "scryfall_id": "12345678-1234-1234-1234-123456789012", "name": "Black Lotus", "card_url": "https://example.com/black-lotus.jpg", "count": 1 }
type Card struct {
	ScryfallID string `json:"scryfall_id" example:"12345678-1234-1234-1234-123456789012"`
	Name       string `json:"name" example:"Black Lotus"`
	CardUrl    string `json:"card_url" example:"https://example.com/black-lotus.jpg"`
	Count      int    `json:"count" example:"1"`
}

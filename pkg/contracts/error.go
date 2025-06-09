package dto

// ErrorResponse — стандартная структура ошибки
// @Description Структура ответа при ошибке
// @example { "message": "unauthorized", "status": 401 }
type ErrorResponse struct {
	Message string `json:"message" example:"unauthorized"`
	Status  int    `json:"status,omitempty"` // Optional, can be used to indicate HTTP status code
}

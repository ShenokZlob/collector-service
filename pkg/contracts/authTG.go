package dto

// RegisterTelegramRequest — данные для регистрации нового пользователя
// @Description Регистрация пользователя по Telegram ID и данным профиля
// @example { "telegram_id": 123456789, "first_name": "Ivan", "last_name": "Ivanov", "username": "ivan123" }
type RegisterTelegramRequest struct {
	TelegramID int64  `json:"telegram_id" binding:"required" example:"123456789"`
	FirstName  string `json:"first_name" binding:"required" example:"Ivan"`
	LastName   string `json:"last_name,omitempty" example:"Ivanov"`
	Username   string `json:"username,omitempty" example:"ivan123"`
}

// RegisterTelegramResponse — ответ после регистрации
// @Description Ответ с JWT-токеном
// @example { "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...", "refresh_token": "dGhpc19pc19hX3JlZnJlc2hfdG9rZW4=" }
type RegisterTelegramResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"dGhpc19pc19hX3JlZnJlc2hfdG9rZW4="`

	// Token string `json:"token" example:"eyJhbG..."`
}

package dto

// RegisterRequest — данные для регистрации нового пользователя по email и паролю
// @Description Регистрация пользователя по email и паролю
// @example { "email": "user@example.com", "password": "strongpassword", "first_name": "Ivan", "last_name": "Ivanov" }
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email" example:"user@example.com"`
	Password  string `json:"password" binding:"required,min=6" example:"strongpassword"`
	FirstName string `json:"first_name" binding:"required" example:"Ivan"`
	LastName  string `json:"last_name,omitempty" example:"Ivanov"`
}

// RegisterResponse — ответ после регистрации
// @Description Ответ с JWT-токенами (access и refresh)
// @example { "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...", "expires_at": "900" }
type RegisterResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt   int    `json:"expires_at" example:"900"`

	// RefreshToken string `json:"refresh_token" example:"dGhpc19pc19hX3JlZnJlc2hfdG9rZW4="`
}

// LoginRequest — данные для входа пользователя по email и паролю
// @Description Вход пользователя по email и паролю
// @example { "email": "user@example.com", "password": "strongpassword" }
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"strongpassword"`
}

// LoginResponse — ответ после успешного входа
// @Description Ответ с JWT-токенами (access и refresh) после логина
// @example { "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...", "expires_at": "900" }
type LoginResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt   int    `json:"expires_at" example:"900"`

	// RefreshToken string `json:"refresh_token" example:"dGhpc19pc19hX3JlZnJlc2hfdG9rZW4="`
}

// RefreshTokenRequest — данные для обновления access-токена
// @Description Обновление access-токена по refresh-токену
// @example { "refresh_token": "dGhpc19pc19hX3JlZnJlc2hfdG9rZW4=" }
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"dGhpc19pc19hX3JlZnJlc2hfdG9rZW4="`
}

// RefreshTokenResponse — ответ с новыми JWT-токенами после обновления
// @Description Новый access- и refresh-токены после успешного refresh
// @example { "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...", "expires_at": "900" }
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt   int    `json:"expires_at" example:"900"`

	// RefreshToken string `json:"refresh_token" example:"bmV3X3JlZnJlc2hfdG9rZW4="`
}

// LogoutRequest — данные для выхода (инвалидации) текущей сессии
// @Description Выход пользователя: инвалидация предоставленного refresh-токена
// @example { "refresh_token": "dGhpc19pc19hX3JlZnJlc2hfdG9rZW4=" }
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"dGhpc19pc19hX3JlZnJlc2hfdG9rZW4="`
}

// LogoutResponse — ответ после успешного выхода
// @Description Подтверждение успешного logout
// @example { "message": "Logout successful" }
type LogoutResponse struct {
	Message string `json:"message" example:"Logout successful"`
}

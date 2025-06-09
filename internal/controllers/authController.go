package controllers

import (
	"net/http"

	"github.com/ShenokZlob/collector-service/domain"
	dto "github.com/ShenokZlob/collector-service/pkg/contracts"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthController отвечает за регистрацию, логин и проверку пользователя
// @Tags Auth
// @BasePath /
type AuthController struct {
	log         *zap.Logger
	authService AuthUsecase
}

type AuthUsecase interface {
	Register(data *dto.RegisterRequest) (string, string, *domain.ResponseErr)
	Login(data *dto.LoginRequest) (string, string, *domain.ResponseErr)
	Refresh(token string) (string, string, *domain.ResponseErr)
	Logout(token string) *domain.ResponseErr

	RegisterTelegram(*domain.User) (string, string, *domain.ResponseErr)
}

func NewAuthController(log *zap.Logger, authService AuthUsecase) *AuthController {
	return &AuthController{
		log:         log.With(zap.String("controller", "auth")),
		authService: authService,
	}
}

// @Summary     Register user by email and password
// @Description Регистрация пользователя по email и паролю
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       input body dto.RegisterRequest true "Данные для регистрации"
// @Success     201 {object} dto.RegisterResponse
// @Failure     400 {object} domain.ResponseErr
// @Router      /register [post]
func (ac AuthController) Register(ctx *gin.Context) {
	ac.log.Info("Register: started")

	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ac.log.Error("Failed to bind json", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, domain.ResponseErr{Message: err.Error()})
		return
	}

	accessToken, refreshToken, respErr := ac.authService.Register(&req)
	if respErr != nil {
		ac.log.Error("Failed to register user", zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ac.log.Info("Register: success", zap.String("accessToken", accessToken),
		zap.String("refreshToken", refreshToken))

	maxAge := 7 * 24 * 3600
	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		maxAge,
		"/",
		"",
		true,
		true,
	)

	ctx.JSON(http.StatusCreated, dto.RegisterResponse{
		AccessToken: accessToken,
		ExpiresAt:   15 * 60,
	})
}

// @Summary     Login user by email
// @Description Вход пользователя по email и паролю
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       input body dto.LoginRequest true "Данные для входа"
// @Success     200 {object} dto.LoginResponse
// @Failure     400 {object} domain.ResponseErr
// @Failure     401 {object} domain.ResponseErr
// @Router      /login [post]
func (ac AuthController) Login(ctx *gin.Context) {
	ac.log.Info("Login: started")

	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ac.log.Error("Failed to bind json", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, domain.ResponseErr{Message: err.Error()})
		return
	}

	accessToken, refreshToken, respErr := ac.authService.Login(&req)
	if respErr != nil {
		ac.log.Error("Failed to login user", zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ac.log.Info("Login: success", zap.String("accessToken", accessToken),
		zap.String("refreshToken", refreshToken))

	maxAge := 7 * 24 * 3600
	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		maxAge,
		"/",
		"",
		true,
		true,
	)

	ctx.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken: accessToken,
		ExpiresAt:   15 * 60,
	})
}

// @Summary     Refresh JWT tokens
// @Description Обновление access и refresh токенов
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       input body dto.RefreshTokenRequest true "Refresh token"
// @Success     200 {object} dto.RefreshTokenResponse
// @Failure     400 {object} domain.ResponseErr
// @Failure     401 {object} domain.ResponseErr
// @Router      /refresh [post]
func (ac AuthController) RefreshToken(ctx *gin.Context) {
	ac.log.Info("RefreshToken: started")

	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ac.log.Error("Failed to bind json", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, domain.ResponseErr{Message: err.Error()})
		return
	}

	accessToken, refreshToken, respErr := ac.authService.Refresh(req.RefreshToken)
	if respErr != nil {
		ac.log.Error("Failed to refresh token", zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ac.log.Info("RefreshToken: success", zap.String("accessToken", accessToken),
		zap.String("refreshToken", refreshToken))

	maxAge := 7 * 24 * 3600
	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		maxAge,
		"/",
		"",
		true,
		true,
	)

	ctx.JSON(http.StatusOK, dto.RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   15 * 60,
	})
}

// @Summary     Logout user
// @Description Выход пользователя и инвалидация refresh токена
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       input body dto.LogoutRequest true "Refresh token для выхода"
// @Success     200 {object} dto.LogoutResponse
// @Failure     400 {object} domain.ResponseErr
// @Failure     401 {object} domain.ResponseErr
// @Router      /logout [post]
func (ac AuthController) Logout(ctx *gin.Context) {
	ac.log.Info("Logout: started")

	var req dto.LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ac.log.Error("Failed to bind json", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, domain.ResponseErr{Message: err.Error()})
		return
	}

	respErr := ac.authService.Logout(req.RefreshToken)
	if respErr != nil {
		ac.log.Error("Failed to logout user", zap.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ac.log.Info("Logout: success")

	ctx.JSON(http.StatusOK, dto.LogoutResponse{
		Message: "Logout successful",
	})
}

// @Summary     Register telegram user
// @Description Регистрация пользователя через Telegram бота
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       input body dto.RegisterRequest true "Данные для регистрации"
// @Success     201 {object} dto.RegisterResponse
// @Failure     400 {object} domain.ResponseErr
// @Router      /register [post]
func (ac AuthController) RegisterTelegram(ctx *gin.Context) {
	ac.log.Info("RegisterTelegram: started")

	var req dto.RegisterTelegramRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ac.log.Error("Failed to bind json", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, domain.ResponseErr{Message: err.Error()})
		return
	}

	userModel := &domain.User{
		TelegramID: req.TelegramID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Username:   req.Username,
	}
	accessToken, refreshToken, respErr := ac.authService.RegisterTelegram(userModel)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ac.log.Info("RegisterTelegram: success", zap.String("accessToken", accessToken),
		zap.String("refreshToken", refreshToken))

	ctx.JSON(http.StatusCreated, dto.RegisterTelegramResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (ac AuthController) LinkTelegram(ctx *gin.Context) {
	panic("TODO")
}

func (ac AuthController) UnlinkTelegram(ctx *gin.Context) {
	panic("TODO")
}

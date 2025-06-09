package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ShenokZlob/collector-service/domain"
	dto "github.com/ShenokZlob/collector-service/pkg/contracts"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	log            *zap.Logger
	authRepository AuthRepositorer
}

type AuthRepositorer interface {
	CreateUser(user *domain.User) (*domain.User, *domain.ResponseErr)
	GetUser(userID string) (*domain.User, *domain.ResponseErr)
	FindByEmail(email string) (*domain.User, *domain.ResponseErr)
	AddToBlackList(jti string) *domain.ResponseErr
	AddToken(userID, jti string, issued, expires time.Time) *domain.ResponseErr

	FindUserByTelegramID(telegramId int64) (*domain.User, *domain.ResponseErr)
}

func NewAuthUsecase(log *zap.Logger, authRepository AuthRepositorer) *AuthUsecase {
	return &AuthUsecase{
		log:            log.With(zap.String("usecase", "auth")),
		authRepository: authRepository,
	}
}

// For simple users

func (as AuthUsecase) Register(data *dto.RegisterRequest) (string, string, *domain.ResponseErr) {
	if !validateRegData(data.Email, data.Password) {
		as.log.Warn("Invalid email or password")
		return "", "", &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid email or password",
		}
	}

	if !validateUser(data.FirstName) {
		as.log.Warn("Invalid user's data")
		return "", "", &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid user's data",
		}
	}

	hash, err := hashPassword(data.Password)
	if err != nil {
		as.log.Error("Failed to hash password", zap.Error(err))
		return "", "", &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: "Failed to hash password",
		}
	}

	user := &domain.User{
		Email:        data.Email,
		PasswordHash: hash,
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, respErr := as.authRepository.CreateUser(user)
	if respErr != nil {
		as.log.Error("Failed to create user", zap.Error(err))
		return "", "", respErr
	}

	issuedAt := time.Now()
	jti := uuid.NewV4().String()
	accessToken, _, err := generateAccessToken(createdUser.ID, jti, issuedAt)
	if err != nil {
		as.log.Error("Failed to generate access token", zap.Error(err))
		return "", "", &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate access token",
		}
	}
	refreshToken, refreshExp, err := generateRefreshToken(createdUser.ID, jti, issuedAt)
	if err != nil {
		as.log.Error("Failed to generate refresh token", zap.Error(err))
		return "", "", &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate refresh token",
		}
	}

	respErr = as.authRepository.AddToken(createdUser.ID, jti, issuedAt, refreshExp)
	if respErr != nil {
		as.log.Error("Failed to add token to db")
		return "", "", respErr
	}

	return accessToken, refreshToken, nil
}

func (as AuthUsecase) Login(data *dto.LoginRequest) (string, string, *domain.ResponseErr) {
	// TODO: add old token to black list
	if !validateRegData(data.Email, data.Password) {
		as.log.Warn("Invalid email or password")
		return "", "", &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid email or password",
		}
	}

	foundUser, respErr := as.authRepository.FindByEmail(data.Email)
	if respErr != nil {
		as.log.Error("Failed to find user by email", zap.String("email", data.Email), zap.Error(respErr))
		return "", "", respErr
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(data.Password)); err != nil {
		as.log.Warn("Passwords hashes are not the same")
		return "", "", &domain.ResponseErr{
			Status:  http.StatusUnauthorized,
			Message: "Invalid email or password",
		}
	}

	issuedAt := time.Now()
	jti := uuid.NewV4().String()
	accessToken, _, err := generateAccessToken(foundUser.ID, jti, issuedAt)
	if err != nil {
		as.log.Error("Failed to generate access token", zap.Error(err))
		return "", "", &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate access token",
		}
	}
	refreshToken, refreshExp, err := generateRefreshToken(foundUser.ID, jti, issuedAt)
	if err != nil {
		as.log.Error("Failed to generate refresh token", zap.Error(err))
		return "", "", &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate refresh token",
		}
	}

	respErr = as.authRepository.AddToken(foundUser.ID, jti, issuedAt, refreshExp)
	if respErr != nil {
		as.log.Error("Failed to add token to db", zap.Error(respErr))
		return "", "", respErr
	}

	return accessToken, refreshToken, nil
}

func (as AuthUsecase) Refresh(token string) (string, string, *domain.ResponseErr) {
	userID, jtiOld, err := checkRefreshToken(token)
	if err != nil {
		as.log.Error("Failed to check refresh token", zap.Error(err))
		return "", "", &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Failed to check refresh token",
		}
	}

	user, respErr := as.authRepository.GetUser(userID)
	if respErr != nil {
		as.log.Error("Failed to find user", zap.Error(respErr))
		return "", "", respErr
	}

	issuedAt := time.Now()
	jtiNew := uuid.NewV4().String()
	accessToken, _, err := generateAccessToken(user.ID, jtiNew, issuedAt)
	if err != nil {
		as.log.Error("Failed to generate access token", zap.Error(err))
		return "", "", &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate access token",
		}
	}
	refreshToken, refreshExp, err := generateRefreshToken(user.ID, jtiNew, issuedAt)
	if err != nil {
		as.log.Error("Failed to generate refresh token", zap.Error(err))
		return "", "", &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate refresh token",
		}
	}

	respErr = as.authRepository.AddToken(userID, jtiNew, issuedAt, refreshExp)
	if respErr != nil {
		as.log.Error("Failed to add new token in db", zap.Error(respErr))
		return "", "", respErr
	}

	respErr = as.authRepository.AddToBlackList(jtiOld)
	if respErr != nil {
		as.log.Error("Failed to add old refresh token to black list", zap.Error(respErr))
		return "", "", respErr
	}

	return accessToken, refreshToken, nil
}

func (as AuthUsecase) Logout(token string) *domain.ResponseErr {
	_, jti, err := checkRefreshToken(token)
	if err != nil {
		as.log.Error("Failed to check refresh token", zap.Error(err))
		return &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Failed to check refresh token",
		}
	}

	respErr := as.authRepository.AddToBlackList(jti)
	if respErr != nil {
		as.log.Error("Failed to add token to black list", zap.Error(respErr))
		return respErr
	}

	return nil
}

// For Telegram users

// Register creates a new user in the database.
func (as AuthUsecase) RegisterTelegram(user *domain.User) (string, string, *domain.ResponseErr) {
	as.log.With(zap.String("method", "Register")).Info("registering user")

	respErr := validateTelegramUser(user)
	if respErr != nil {
		as.log.Error("failed to validate user", zap.Error(respErr))
		return "", "", respErr
	}

	createdUser, respErr := as.authRepository.CreateUser(user)
	if respErr != nil {
		as.log.Error("failed to create user", zap.Error(respErr))
		return "", "", respErr
	}

	issuedAt := time.Now()
	jtiNew := uuid.NewV4().String()
	accessToken, _, err := generateAccessToken(createdUser.ID, jtiNew, issuedAt)
	if err != nil {
		as.log.Error("Failed to generate access token", zap.Error(err))
		return "", "", &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate access token",
		}
	}
	refreshToken, refreshExp, err := generateRefreshToken(createdUser.ID, jtiNew, issuedAt)
	if err != nil {
		as.log.Error("Failed to generate refresh token", zap.Error(err))
		return "", "", &domain.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate refresh token",
		}
	}

	respErr = as.authRepository.AddToken(createdUser.ID, jtiNew, issuedAt, refreshExp)
	if respErr != nil {
		as.log.Error("Failed to add token to db")
		return "", "", respErr
	}

	return accessToken, refreshToken, nil
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func generateAccessToken(userID, jti string, issuedAt time.Time) (string, time.Time, error) {
	expiresAt := issuedAt.Add(15 * time.Minute)
	token, err := generateJWTToken(userID, jti, "", issuedAt, expiresAt)
	return token, expiresAt, err
}

func generateRefreshToken(userID, jti string, issuedAt time.Time) (string, time.Time, error) {
	expiresAt := issuedAt.Add(7 * 24 * time.Hour)
	token, err := generateJWTToken(userID, jti, "refresh", issuedAt, expiresAt)
	return token, expiresAt, err
}

func generateJWTToken(userID, jti, tokenType string, issuedAt, expiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"jti": jti,
		"iss": issuedAt.Unix(),
		"exp": expiresAt.Unix(),
	}
	if tokenType != "" {
		claims["type"] = tokenType
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

// checkRefreshToken checks token and return sub
func checkRefreshToken(tokenStr string) (userID, jti string, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid JWT claims")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", "", fmt.Errorf("not refresh token")
	}

	userID, ok = claims["sub"].(string)
	if !ok || userID == "" {
		return "", "", fmt.Errorf("invalid user ID")
	}

	jti, ok = claims["jti"].(string)
	if !ok || jti == "" {
		return "", "", fmt.Errorf("invalid jti")
	}

	return
}

func validateRegData(email, password string) bool {
	// TODO: replace
	return email != "" && password != ""
}

func validateUser(firstName string) bool {
	return firstName != ""
}

func validateTelegramUser(user *domain.User) *domain.ResponseErr {
	if user.TelegramID == 0 {
		return &domain.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid user telegram ID",
		}
	}

	return nil
}

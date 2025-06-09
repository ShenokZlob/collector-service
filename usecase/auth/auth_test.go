package auth

import (
	"testing"
	"time"

	"github.com/ShenokZlob/collector-service/domain"
	dto "github.com/ShenokZlob/collector-service/pkg/contracts"
	"github.com/ShenokZlob/collector-service/usecase/auth/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

func TestGenerateAccessToken(t *testing.T) {
	userID := "userID"
	jti := "jtiTest"
	issuedAt := time.Now()

	// With JWT_SECRET
	t.Setenv("JWT_SECRET", "testsecret")
	token, expAt, err := generateAccessToken(userID, jti, issuedAt)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.WithinDuration(t, issuedAt.Add(15*time.Minute), expAt, time.Second)

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("testsecret"), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsed.Valid)

	claims, ok := parsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, userID, claims["sub"])
	assert.Equal(t, jti, claims["jti"])
	assert.WithinDuration(t, issuedAt, time.Unix(int64(claims["iss"].(float64)), 0), time.Second)
	assert.WithinDuration(t, expAt, time.Unix(int64(claims["exp"].(float64)), 0), time.Second)
	assert.Nil(t, claims["type"])
}

func TestGenerateRefreshToken(t *testing.T) {
	userID := "userID"
	jti := "jtiTest"
	issuedAt := time.Now()

	// With JWT_SECRET
	t.Setenv("JWT_SECRET", "testsecret")
	token, expAt, err := generateRefreshToken(userID, jti, issuedAt)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.WithinDuration(t, issuedAt.Add(7*24*time.Hour), expAt, time.Second)

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("testsecret"), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsed.Valid)

	claims, ok := parsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, userID, claims["sub"])
	assert.Equal(t, jti, claims["jti"])
	assert.WithinDuration(t, issuedAt, time.Unix(int64(claims["iss"].(float64)), 0), time.Second)
	assert.WithinDuration(t, expAt, time.Unix(int64(claims["exp"].(float64)), 0), time.Second)
	assert.NotNil(t, claims["type"])
}

type UnitySuite struct {
	suite.Suite
}

func TestUnitySuite(t *testing.T) {
	suite.Run(t, &UnitySuite{})
}

func TestRegister(t *testing.T) {
	// Setup
	t.Setenv("JWT_SECRET", "testsecret")
	var repMock mocks.MockAuthRepositorer

	expectedUser := &domain.User{
		ID:           "user123",
		Email:        "email@test.com",
		PasswordHash: "",
		FirstName:    "Test",
		LastName:     "Testovic",
	}

	repMock.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(expectedUser, nil)
	repMock.On("AddToken", expectedUser.ID, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	authUsecase := NewAuthUsecase(zap.NewNop(), &repMock)

	dataTest := &dto.RegisterRequest{
		Email:     "email@test.com",
		Password:  "testpassword",
		FirstName: "Test",
		LastName:  "Testovic",
	}

	accessToken, refreshToken, err := authUsecase.Register(dataTest)

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	repMock.AssertExpectations(t)
}

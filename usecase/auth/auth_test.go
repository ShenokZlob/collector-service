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
	"golang.org/x/crypto/bcrypt"
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
	au      *AuthUsecase
	repMock *mocks.MockAuthRepositorer
}

func TestUnitySuite(t *testing.T) {
	suite.Run(t, &UnitySuite{})
}

func (us *UnitySuite) SetupSuite() {
	us.T().Setenv("JWT_SECRET", "testsecret")
}

func (us *UnitySuite) SetupTest() {
	us.repMock = &mocks.MockAuthRepositorer{}
	us.au = NewAuthUsecase(zap.NewNop(), us.repMock)
}

func (us *UnitySuite) TestRegister() {
	expectedUser := &domain.User{
		ID:           "user123",
		Email:        "email@test.com",
		PasswordHash: "",
		FirstName:    "Test",
		LastName:     "Testovic",
	}

	us.repMock.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(expectedUser, nil)
	us.repMock.On("AddToken", expectedUser.ID, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	dataTest := &dto.RegisterRequest{
		Email:     "email@test.com",
		Password:  "testpassword",
		FirstName: "Test",
		LastName:  "Testovic",
	}

	accessToken, refreshToken, respErr := us.au.Register(dataTest)

	us.Assert().Nil(respErr)
	us.Assert().NotEmpty(accessToken)
	us.Assert().NotEmpty(refreshToken)
	us.repMock.AssertExpectations(us.T())
}

func (us *UnitySuite) TestLogin() {
	password := "testpassword"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	expectedUser := &domain.User{
		ID:           "user123",
		Email:        "email@test.com",
		PasswordHash: string(hash),
		FirstName:    "Test",
		LastName:     "Testovic",
	}

	us.repMock.On("FindByEmail", mock.Anything).Return(expectedUser, nil)
	us.repMock.On("AddToken", expectedUser.ID, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	data := &dto.LoginRequest{
		Email:    "email@test.com",
		Password: password,
	}

	accessToken, refreshToken, respErr := us.au.Login(data)

	us.Assert().Nil(respErr)
	us.Assert().NotEmpty(accessToken)
	us.Assert().NotEmpty(refreshToken)
	us.repMock.AssertExpectations(us.T())
}

func (us *UnitySuite) TestRefresh() {
	expectedUser := &domain.User{
		ID:           "user123",
		Email:        "email@test.com",
		PasswordHash: "",
		FirstName:    "Test",
		LastName:     "Testovic",
	}

	issAt := time.Now()
	refreshToken, expAt, err := generateRefreshToken(expectedUser.ID, "jtitest", issAt)
	us.Assert().Equal(expAt, issAt.Add(7*24*time.Hour))
	us.Require().NoError(err)

	us.repMock.On("GetUser", expectedUser.ID).Return(expectedUser, nil)
	us.repMock.On("AddToken", expectedUser.ID, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	us.repMock.On("AddToBlackList", mock.Anything).Return(nil)

	accessToken, refreshToken, respErr := us.au.Refresh(refreshToken)

	us.Assert().Nil(respErr)
	us.Assert().NotEmpty(accessToken)
	us.Assert().NotEmpty(refreshToken)
	us.repMock.AssertExpectations(us.T())
}

func (us *UnitySuite) TestLogout() {
	expectedUser := &domain.User{
		ID:           "user123",
		Email:        "email@test.com",
		PasswordHash: "",
		FirstName:    "Test",
		LastName:     "Testovic",
	}

	issAt := time.Now()
	refreshToken, expAt, err := generateRefreshToken(expectedUser.ID, "jtitest", issAt)
	us.Assert().Equal(expAt, issAt.Add(7*24*time.Hour))
	us.Require().NoError(err)

	us.repMock.On("AddToBlackList", mock.Anything).Return(nil)

	respErr := us.au.Logout(refreshToken)

	us.Assert().Nil(respErr)
	us.repMock.AssertExpectations(us.T())
}

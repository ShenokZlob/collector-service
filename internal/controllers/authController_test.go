package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mocks "github.com/ShenokZlob/collector-service/internal/controllers/mocks"
	"github.com/gin-gonic/gin"
	faker "github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRegister(t *testing.T) {
	// Arrange
	mockAuthService := new(mocks.MockAuthUsecase)
	ctrl := AuthController{
		log:         zap.NewNop(),
		authService: mockAuthService,
	}

	email := faker.Email()
	password := faker.Password()
	firstName := faker.FirstName()
	lastName := faker.LastName()

	reqBody := `{"email":"` + email + `","password":"` + password + `","first_name":"` + firstName + `","last_name":"` + lastName + `"}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	accessToken := "access.jwt.token"
	refreshToken := "refresh.jwt.token"

	mockAuthService.
		On("Register", mock.AnythingOfType("*dto.RegisterRequest")).
		Return(accessToken, refreshToken, nil)

	// Act
	ctrl.Register(c)

	// Assert
	require.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), accessToken)
	//assert.Contains(t, w.Body.String(), refreshToken) // in header
	mockAuthService.AssertExpectations(t)
}

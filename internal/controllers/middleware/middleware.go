package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func AuthMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Info("JWT middleware triggered")

		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Info("Authorization header does not start with Bearer", zap.String("header", authHeader))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			log.Error("Failed to parse JWT token", zap.Error(err), zap.String("token", tokenStr))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Error("Invalid JWT claims", zap.String("claims", tokenStr))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			log.Error("Invalid or missing user_id in JWT claims")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
			log.Warn("JWT token has expired", zap.Int("expiration", int(exp)), zap.String("current_time", time.Now().String()))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}

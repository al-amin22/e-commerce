package middleware

import (
	"net/http"
	"strings"

	"ecommerce-backend/pkg/security"
	"github.com/gin-gonic/gin"
)

type JWTAuthMiddleware struct {
	secret string
}

func NewJWTAuthMiddleware(secret string) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{secret: secret}
}

func (m *JWTAuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "missing or invalid authorization header"})
			c.Abort()
			return
		}

		claims, err := security.ParseToken(parts[1], m.secret)
		if err != nil || claims.Type != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired access token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

package middleware

import (
	"net/http"
	"strings"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: secret}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := ""
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenStr = parts[1]
			}
		}

		if tokenStr == "" {
			if cookieToken, err := c.Cookie("access_token"); err == nil {
				tokenStr = cookieToken
			}
		}

		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "missing access token"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenStr, m.jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
			c.Abort()
			return
		}

		if claims.Type != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token type"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...models.UserRole) gin.HandlerFunc {
	allowed := make(map[models.UserRole]bool)
	for _, role := range roles {
		allowed[role] = true
	}

	return func(c *gin.Context) {
		roleAny, ok := c.Get("role")
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"message": "role not found"})
			c.Abort()
			return
		}

		role, ok := roleAny.(models.UserRole)
		if !ok || !allowed[role] {
			c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}

package utils

import (
	"time"

	"ecommerce-backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID uuid.UUID       `json:"user_id"`
	Role   models.UserRole `json:"role"`
	Type   string          `json:"type"`
	JTI    string          `json:"jti,omitempty"`
	jwt.RegisteredClaims
}

func CreateAccessToken(userID uuid.UUID, role models.UserRole, secret string, ttl time.Duration) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func CreateRefreshToken(userID uuid.UUID, role models.UserRole, secret string, ttl time.Duration) (string, string, error) {
	jti := uuid.NewString()
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		Type:   "refresh",
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return signed, jti, nil
}

func ParseToken(tokenStr, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

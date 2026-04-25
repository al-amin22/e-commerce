package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/repository"
	"ecommerce-backend/pkg/config"
	"ecommerce-backend/pkg/security"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, *AuthTokens, error)
	Refresh(ctx context.Context, refreshToken string) (*AuthTokens, error)
	GetUserByID(ctx context.Context, id uint) (*domain.User, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	cfg      *config.Config
	userRepo repository.UserRepository
	redis    *redis.Client
}

func NewAuthService(cfg *config.Config, userRepo repository.UserRepository, redisClient *redis.Client) AuthService {
	return &authService{cfg: cfg, userRepo: userRepo, redis: redisClient}
}

func (s *authService) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || strings.TrimSpace(password) == "" || strings.TrimSpace(name) == "" {
		return nil, errors.New("name, email, and password are required")
	}

	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:         strings.TrimSpace(name),
		Email:        email,
		PasswordHash: string(hashed),
		Role:         "buyer",
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*domain.User, *AuthTokens, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	tokens, err := s.issueTokenPair(ctx, user.ID, user.Role)
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	claims, err := security.ParseToken(refreshToken, s.cfg.JWTSecret)
	if err != nil || claims.Type != "refresh" || claims.ID == "" {
		return nil, errors.New("invalid refresh token")
	}

	activeJTI, err := s.redis.Get(ctx, s.activeRefreshKey(claims.UserID)).Result()
	if err != nil || activeJTI != claims.ID {
		return nil, errors.New("refresh token revoked")
	}

	if err := s.redis.Del(ctx, s.refreshKey(claims.UserID, claims.ID)).Err(); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.issueTokenPair(ctx, user.ID, user.Role)
}

func (s *authService) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := security.ParseToken(refreshToken, s.cfg.JWTSecret)
	if err != nil || claims.Type != "refresh" || claims.ID == "" {
		return nil
	}

	if err := s.redis.Del(ctx, s.refreshKey(claims.UserID, claims.ID)).Err(); err != nil {
		return err
	}
	if err := s.redis.Del(ctx, s.activeRefreshKey(claims.UserID)).Err(); err != nil {
		return err
	}
	return nil
}

func (s *authService) issueTokenPair(ctx context.Context, userID uint, role string) (*AuthTokens, error) {
	accessTTL := time.Duration(s.cfg.AccessTokenTTLMin) * time.Minute
	refreshTTL := time.Duration(s.cfg.RefreshTokenTTLDays) * 24 * time.Hour
	jti := uuid.NewString()

	accessToken, err := security.GenerateToken(userID, role, "access", s.cfg.JWTSecret, accessTTL, "")
	if err != nil {
		return nil, err
	}

	refreshToken, err := security.GenerateToken(userID, role, "refresh", s.cfg.JWTSecret, refreshTTL, jti)
	if err != nil {
		return nil, err
	}

	if prevJTI, err := s.redis.Get(ctx, s.activeRefreshKey(userID)).Result(); err == nil && prevJTI != "" {
		s.redis.Del(ctx, s.refreshKey(userID, prevJTI))
	}

	if err := s.redis.Set(ctx, s.refreshKey(userID, jti), "1", refreshTTL).Err(); err != nil {
		return nil, err
	}
	if err := s.redis.Set(ctx, s.activeRefreshKey(userID), jti, refreshTTL).Err(); err != nil {
		return nil, err
	}

	return &AuthTokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *authService) refreshKey(userID uint, jti string) string {
	return fmt.Sprintf("auth:refresh:%d:%s", userID, jti)
}

func (s *authService) activeRefreshKey(userID uint) string {
	return fmt.Sprintf("auth:refresh:active:%d", userID)
}

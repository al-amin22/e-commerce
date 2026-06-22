package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"ecommerce-backend/internal/config"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/utils"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthService struct {
	db    *gorm.DB
	redis *redis.Client
	cfg   *config.Config
}

func NewAuthService(db *gorm.DB, redisClient *redis.Client, cfg *config.Config) *AuthService {
	return &AuthService{db: db, redis: redisClient, cfg: cfg}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*models.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)
	if name == "" || email == "" || strings.TrimSpace(password) == "" {
		return nil, "", errors.New("name, email, and password are required")
	}

	if _, err := s.findUserByEmail(ctx, email); err == nil {
		return nil, "", errors.New("email already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         models.RoleBuyer,
		IsVerified:   false,
	}

	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, "", err
	}

	otp := &models.EmailOTP{
		Email:     email,
		Code:      utils.GenerateOTP(6),
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := s.db.WithContext(ctx).Create(otp).Error; err != nil {
		return nil, "", err
	}

	return user, otp.Code, nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, email, otpCode string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || strings.TrimSpace(otpCode) == "" {
		return errors.New("email and otp are required")
	}

	var otp models.EmailOTP
	err := s.db.WithContext(ctx).
		Where("email = ? AND code = ? AND used = ?", email, otpCode, false).
		Order("created_at DESC").
		First(&otp).Error
	if err != nil {
		return errors.New("invalid otp")
	}

	if otp.ExpiresAt.Before(time.Now()) {
		return errors.New("otp expired")
	}

	if err := s.db.WithContext(ctx).Model(&otp).Update("used", true).Error; err != nil {
		return err
	}

	if err := s.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Update("is_verified", true).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password, deviceID string) (*models.User, string, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || strings.TrimSpace(password) == "" {
		return nil, "", "", errors.New("email and password are required")
	}

	user, err := s.findUserByEmail(ctx, email)
	if err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	if err := utils.CheckPassword(user.PasswordHash, password); err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return nil, "", "", errors.New("email not verified")
	}

	accessToken, refreshToken, _, err := s.issueTokenPair(ctx, user.ID, user.Role, deviceID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken, deviceID string) (string, string, error) {
	claims, err := utils.ParseToken(refreshToken, s.cfg.JWTSecret)
	if err != nil || claims.Type != "refresh" || claims.JTI == "" {
		return "", "", errors.New("invalid refresh token")
	}

	if ok := s.isRefreshTokenActive(ctx, claims.UserID.String(), claims.JTI, deviceID); !ok {
		return "", "", errors.New("refresh token revoked")
	}

	s.revokeRefreshToken(ctx, claims.UserID.String(), claims.JTI, deviceID)

	user, err := s.findUserByID(ctx, claims.UserID)
	if err != nil {
		return "", "", err
	}

	accessToken, newRefreshToken, _, err := s.issueTokenPair(ctx, user.ID, user.Role, deviceID)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken, deviceID string) error {
	if refreshToken == "" {
		return nil
	}

	claims, err := utils.ParseToken(refreshToken, s.cfg.JWTSecret)
	if err != nil || claims.Type != "refresh" || claims.JTI == "" {
		return nil
	}

	s.revokeRefreshToken(ctx, claims.UserID.String(), claims.JTI, deviceID)
	return nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.findUserByID(ctx, id)
}

func (s *AuthService) issueTokenPair(ctx context.Context, userID uuid.UUID, role models.UserRole, deviceID string) (string, string, string, error) {
	accessTTL := time.Duration(s.cfg.AccessTokenTTLMin) * time.Minute
	refreshTTL := time.Duration(s.cfg.RefreshTokenTTLDays) * 24 * time.Hour

	accessToken, err := utils.CreateAccessToken(userID, role, s.cfg.JWTSecret, accessTTL)
	if err != nil {
		return "", "", "", err
	}

	refreshToken, jti, err := utils.CreateRefreshToken(userID, role, s.cfg.JWTSecret, refreshTTL)
	if err != nil {
		return "", "", "", err
	}

	if prevJTI, err := s.redis.Get(ctx, refreshDeviceKey(userID.String(), deviceID)).Result(); err == nil && prevJTI != "" {
		s.redis.Del(ctx, refreshTokenKey(userID.String(), prevJTI))
	}

	if err := s.redis.Set(ctx, refreshTokenKey(userID.String(), jti), deviceID, refreshTTL).Err(); err != nil {
		return "", "", "", err
	}

	if err := s.redis.Set(ctx, refreshDeviceKey(userID.String(), deviceID), jti, refreshTTL).Err(); err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, jti, nil
}

func (s *AuthService) isRefreshTokenActive(ctx context.Context, userID, jti, deviceID string) bool {
	storedDeviceID, err := s.redis.Get(ctx, refreshTokenKey(userID, jti)).Result()
	if err != nil {
		return false
	}
	return storedDeviceID == deviceID
}

func (s *AuthService) revokeRefreshToken(ctx context.Context, userID, jti, deviceID string) {
	s.redis.Del(ctx, refreshTokenKey(userID, jti))
	s.redis.Del(ctx, refreshDeviceKey(userID, deviceID))
}

func (s *AuthService) findUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) findUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func refreshTokenKey(userID, jti string) string {
	return fmt.Sprintf("refresh:%s:%s", userID, jti)
}

func refreshDeviceKey(userID, deviceID string) string {
	return fmt.Sprintf("refresh_device:%s:%s", userID, deviceID)
}

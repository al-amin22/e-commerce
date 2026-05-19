package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type registerRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type verifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type authResponse struct {
	User gin.H `json:"user"`
}

func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed hashing password"})
		return
	}

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         models.RoleBuyer,
		IsVerified:   false,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "email already registered"})
		return
	}

	otype := models.EmailOTP{
		Email:     req.Email,
		Code:      utils.GenerateOTP(6),
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := h.DB.Create(&otype).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed creating otp"})
		return
	}

	h.EmailSvc.SendOTPAsync(req.Email, otype.Code)

	c.JSON(http.StatusCreated, gin.H{
		"message": "registration successful, check your email for otp",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	var req verifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var otp models.EmailOTP
	err := h.DB.Where("email = ? AND code = ? AND used = ?", req.Email, req.OTP, false).
		Order("created_at DESC").
		First(&otp).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid otp"})
		return
	}

	if otp.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "otp expired"})
		return
	}

	if err := h.DB.Model(&otp).Update("used", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed updating otp"})
		return
	}

	if err := h.DB.Model(&models.User{}).Where("email = ?", req.Email).Update("is_verified", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed verifying user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusForbidden, gin.H{"message": "email not verified"})
		return
	}

	deviceID := resolveDeviceID(c)
	accessToken, refreshToken, _, err := h.issueTokenPair(c.Request.Context(), user.ID, user.Role, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed generating token pair"})
		return
	}

	h.setAuthCookies(c, accessToken, refreshToken, deviceID)

	c.JSON(http.StatusOK, gin.H{
		"message":      "login success",
		"access_token": accessToken,
		"user":         userPayload(user),
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing refresh token"})
		return
	}

	deviceID := resolveDeviceID(c)
	if cookieDeviceID, err := c.Cookie("device_id"); err == nil && cookieDeviceID != "" {
		deviceID = cookieDeviceID
	}

	claims, err := utils.ParseToken(refreshToken, h.Config.JWTSecret)
	if err != nil || claims.Type != "refresh" || claims.JTI == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid refresh token"})
		return
	}

	if ok := h.isRefreshTokenActive(c.Request.Context(), claims.UserID.String(), claims.JTI, deviceID); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "refresh token revoked"})
		return
	}

	h.revokeRefreshToken(c.Request.Context(), claims.UserID.String(), claims.JTI, deviceID)

	accessToken, newRefreshToken, _, err := h.issueTokenPair(c.Request.Context(), claims.UserID, claims.Role, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed rotating token"})
		return
	}

	h.setAuthCookies(c, accessToken, newRefreshToken, deviceID)
	c.JSON(http.StatusOK, gin.H{"message": "token refreshed"})
}

func (h *Handler) Logout(c *gin.Context) {
	deviceID := resolveDeviceID(c)
	if cookieDeviceID, err := c.Cookie("device_id"); err == nil && cookieDeviceID != "" {
		deviceID = cookieDeviceID
	}

	if refreshToken, err := c.Cookie("refresh_token"); err == nil && refreshToken != "" {
		if claims, parseErr := utils.ParseToken(refreshToken, h.Config.JWTSecret); parseErr == nil && claims.JTI != "" {
			h.revokeRefreshToken(c.Request.Context(), claims.UserID.String(), claims.JTI, deviceID)
		}
	}

	h.clearAuthCookies(c)
	c.JSON(http.StatusOK, gin.H{"message": "logout success"})
}

func (h *Handler) Me(c *gin.Context) {
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uuid.UUID)

	var user models.User
	if err := h.DB.Select("id", "name", "email", "role", "is_verified", "created_at", "updated_at").Where("id = ?", userID.String()).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userPayload(user)})
}

func userPayload(user models.User) gin.H {
	return gin.H{
		"id":          user.ID,
		"name":        user.Name,
		"email":       user.Email,
		"role":        user.Role,
		"is_verified": user.IsVerified,
	}
}

func resolveDeviceID(c *gin.Context) string {
	deviceID := strings.TrimSpace(c.GetHeader("X-Device-ID"))
	if deviceID != "" {
		return deviceID
	}
	ua := strings.TrimSpace(c.GetHeader("User-Agent"))
	if ua == "" {
		return "unknown-device"
	}
	if len(ua) > 80 {
		return ua[:80]
	}
	return ua
}

func refreshTokenKey(userID, jti string) string {
	return fmt.Sprintf("refresh:%s:%s", userID, jti)
}

func refreshDeviceKey(userID, deviceID string) string {
	return fmt.Sprintf("refresh_device:%s:%s", userID, deviceID)
}

func (h *Handler) issueTokenPair(ctx context.Context, userID uuid.UUID, role models.UserRole, deviceID string) (string, string, string, error) {
	accessTTL := time.Duration(h.Config.AccessTokenTTLMin) * time.Minute
	refreshTTL := time.Duration(h.Config.RefreshTokenTTLDays) * 24 * time.Hour

	accessToken, err := utils.CreateAccessToken(userID, role, h.Config.JWTSecret, accessTTL)
	if err != nil {
		return "", "", "", err
	}

	refreshToken, jti, err := utils.CreateRefreshToken(userID, role, h.Config.JWTSecret, refreshTTL)
	if err != nil {
		return "", "", "", err
	}

	if prevJTI, err := h.Redis.Get(ctx, refreshDeviceKey(userID.String(), deviceID)).Result(); err == nil && prevJTI != "" {
		h.Redis.Del(ctx, refreshTokenKey(userID.String(), prevJTI))
	}

	if err := h.Redis.Set(ctx, refreshTokenKey(userID.String(), jti), deviceID, refreshTTL).Err(); err != nil {
		return "", "", "", err
	}

	if err := h.Redis.Set(ctx, refreshDeviceKey(userID.String(), deviceID), jti, refreshTTL).Err(); err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, jti, nil
}

func (h *Handler) isRefreshTokenActive(ctx context.Context, userID, jti, deviceID string) bool {
	storedDeviceID, err := h.Redis.Get(ctx, refreshTokenKey(userID, jti)).Result()
	if err != nil {
		return false
	}
	return storedDeviceID == deviceID
}

func (h *Handler) revokeRefreshToken(ctx context.Context, userID, jti, deviceID string) {
	h.Redis.Del(ctx, refreshTokenKey(userID, jti))
	h.Redis.Del(ctx, refreshDeviceKey(userID, deviceID))
}

func (h *Handler) setAuthCookies(c *gin.Context, accessToken, refreshToken, deviceID string) {
	accessMaxAge := h.Config.AccessTokenTTLMin * 60
	refreshMaxAge := h.Config.RefreshTokenTTLDays * 24 * 60 * 60

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.Config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   accessMaxAge,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.Config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   refreshMaxAge,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "device_id",
		Value:    deviceID,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.Config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   refreshMaxAge,
	})
}

func (h *Handler) clearAuthCookies(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{Name: "access_token", Value: "", Path: "/", HttpOnly: true, Secure: h.Config.CookieSecure, SameSite: http.SameSiteLaxMode, MaxAge: -1})
	http.SetCookie(c.Writer, &http.Cookie{Name: "refresh_token", Value: "", Path: "/", HttpOnly: true, Secure: h.Config.CookieSecure, SameSite: http.SameSiteLaxMode, MaxAge: -1})
	http.SetCookie(c.Writer, &http.Cookie{Name: "device_id", Value: "", Path: "/", HttpOnly: true, Secure: h.Config.CookieSecure, SameSite: http.SameSiteLaxMode, MaxAge: -1})
}


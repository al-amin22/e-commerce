package handlers

import (
	"net/http"
	"strings"

	"ecommerce-backend/internal/dto"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, otpCode, err := h.AuthSvc.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "email already registered" {
			status = http.StatusConflict
		}
		response.Error(c, status, err.Error())
		return
	}

	h.EmailSvc.SendOTPAsync(req.Email, otpCode)
	response.Success(c, http.StatusCreated, "registration successful, check your email for otp", gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.AuthSvc.VerifyEmail(c.Request.Context(), req.Email, req.OTP); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "email verified successfully", nil)
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	deviceID := resolveDeviceID(c)
	user, accessToken, refreshToken, err := h.AuthSvc.Login(c.Request.Context(), req.Email, req.Password, deviceID)
	if err != nil {
		status := http.StatusUnauthorized
		if err.Error() == "email not verified" {
			status = http.StatusForbidden
		}
		response.Error(c, status, err.Error())
		return
	}

	h.setAuthCookies(c, accessToken, refreshToken, deviceID)
	response.Success(c, http.StatusOK, "login success", gin.H{
		"access_token": accessToken,
		"user":         userPayload(*user),
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		response.Error(c, http.StatusUnauthorized, "missing refresh token")
		return
	}

	deviceID := resolveDeviceID(c)
	if cookieDeviceID, err := c.Cookie("device_id"); err == nil && cookieDeviceID != "" {
		deviceID = cookieDeviceID
	}

	accessToken, newRefreshToken, err := h.AuthSvc.Refresh(c.Request.Context(), refreshToken, deviceID)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.setAuthCookies(c, accessToken, newRefreshToken, deviceID)
	response.Success(c, http.StatusOK, "token refreshed", nil)
}

func (h *Handler) Logout(c *gin.Context) {
	deviceID := resolveDeviceID(c)
	if cookieDeviceID, err := c.Cookie("device_id"); err == nil && cookieDeviceID != "" {
		deviceID = cookieDeviceID
	}

	if refreshToken, err := c.Cookie("refresh_token"); err == nil && refreshToken != "" {
		_ = h.AuthSvc.Logout(c.Request.Context(), refreshToken, deviceID)
	}

	h.clearAuthCookies(c)
	response.Success(c, http.StatusOK, "logout success", nil)
}

func (h *Handler) Me(c *gin.Context) {
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uuid.UUID)

	user, err := h.AuthSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user not found")
		return
	}

	response.Success(c, http.StatusOK, "success", gin.H{"user": userPayload(*user)})
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


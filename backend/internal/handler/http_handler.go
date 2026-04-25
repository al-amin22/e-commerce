package handler

import (
	"net/http"
	"strings"

	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	authService    service.AuthService
	userService    service.UserService
	productService service.ProductService
	cookieSecure   bool
}

func NewHTTPHandler(authService service.AuthService, userService service.UserService, productService service.ProductService, cookieSecure bool) *HTTPHandler {
	return &HTTPHandler{authService: authService, userService: userService, productService: productService, cookieSecure: cookieSecure}
}

func (h *HTTPHandler) RegisterRoutes(r *gin.Engine, authMW *middleware.JWTAuthMiddleware) {
	r.GET("/health", h.health)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.register)
			auth.POST("/login", h.login)
			auth.POST("/refresh", h.refresh)
			auth.POST("/logout", h.logout)
		}

		api.GET("/products", h.listProducts)

		protected := api.Group("")
		protected.Use(authMW.RequireAuth())
		{
			protected.GET("/auth/me", h.me)
			protected.GET("/users", h.listUsers)
			protected.POST("/products", h.createProduct)
		}
	}
}

func (h *HTTPHandler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type registerRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *HTTPHandler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Name, strings.ToLower(strings.TrimSpace(req.Email)), req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "register success", "user": user})
}

func (h *HTTPHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, tokens, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	h.setRefreshCookie(c, tokens.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"access_token": tokens.AccessToken, "user": user})
}

func (h *HTTPHandler) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "missing refresh token"})
		return
	}

	tokens, err := h.authService.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	h.setRefreshCookie(c, tokens.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"access_token": tokens.AccessToken})
}

func (h *HTTPHandler) logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	_ = h.authService.Logout(c.Request.Context(), refreshToken)
	h.clearRefreshCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logout success"})
}

func (h *HTTPHandler) me(c *gin.Context) {
	userIDAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	userID, ok := userIDAny.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user context"})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *HTTPHandler) listUsers(c *gin.Context) {
	users, err := h.userService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *HTTPHandler) setRefreshCookie(c *gin.Context, token string) {
	c.SetCookie("refresh_token", token, 7*24*60*60, "/", "", h.cookieSecure, true)
}

func (h *HTTPHandler) clearRefreshCookie(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "", h.cookieSecure, true)
}

func (h *HTTPHandler) createProduct(c *gin.Context) {
	var req domain.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.productService.Create(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "product created", "data": req})
}

func (h *HTTPHandler) listProducts(c *gin.Context) {
	products, err := h.productService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

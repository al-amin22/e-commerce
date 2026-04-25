package main

import (
	"log"
	"net/http"

	"ecommerce-backend/internal/cache"
	"ecommerce-backend/internal/config"
	"ecommerce-backend/internal/db"
	"ecommerce-backend/internal/handlers"
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database, err := db.Connect(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	redisClient, err := cache.ConnectRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("redis error: %v", err)
	}

	emailSvc := services.NewEmailService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.FromEmail)
	shippingSvc := services.NewShippingService(cfg.RajaOngkirAPIKey, cfg.RajaOngkirBase)
	h := handlers.NewHandler(database, redisClient, cfg, emailSvc, shippingSvc)
	authMW := middleware.NewAuthMiddleware(cfg.JWTSecret)

	r := gin.Default()
	r.Use(corsMiddleware(cfg.FrontendOrigin))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/verify-email", h.VerifyEmail)
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.Refresh)
			auth.POST("/logout", h.Logout)
		}

		api.POST("/bootstrap-admin", h.BootstrapAdmin)

		secured := api.Group("")
		secured.Use(authMW.RequireAuth())
		{
			secured.GET("/auth/me", h.Me)
			secured.GET("/addresses", h.ListMyAddresses)
			secured.POST("/addresses", h.CreateAddress)
			secured.GET("/shipping/cost", h.GetShippingCost)

			admin := secured.Group("/admin")
			admin.Use(middleware.RequireRole(models.RoleAdmin))
			admin.GET("/users", h.ListUsers)
		}
	}

	log.Printf("server running at :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(frontendOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == frontendOrigin {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Device-ID")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

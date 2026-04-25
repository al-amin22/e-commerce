package main

import (
	"log"
	"net/http"

	"ecommerce-backend/internal/handler"
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/repository"
	"ecommerce-backend/internal/service"
	"ecommerce-backend/pkg/cache"
	"ecommerce-backend/pkg/config"
	"ecommerce-backend/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := database.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	redisClient, err := cache.ConnectRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("redis error: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)

	authService := service.NewAuthService(cfg, userRepo, redisClient)
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	authMW := middleware.NewJWTAuthMiddleware(cfg.JWTSecret)

	h := handler.NewHTTPHandler(authService, userService, productService, cfg.CookieSecure)

	r := gin.Default()
	r.Use(corsMiddleware(cfg.FrontendOrigin))
	h.RegisterRoutes(r, authMW)

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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

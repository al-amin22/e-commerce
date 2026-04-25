package handlers

import (
	"ecommerce-backend/internal/config"
	"ecommerce-backend/internal/services"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	DB       *gorm.DB
	Redis    *redis.Client
	Config   *config.Config
	EmailSvc *services.EmailService
	ShipSvc  *services.ShippingService
}

func NewHandler(db *gorm.DB, redisClient *redis.Client, cfg *config.Config, email *services.EmailService, ship *services.ShippingService) *Handler {
	return &Handler{DB: db, Redis: redisClient, Config: cfg, EmailSvc: email, ShipSvc: ship}
}

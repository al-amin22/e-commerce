package db

import (
	"fmt"

	"ecommerce-backend/internal/config"
	"ecommerce-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	database, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := database.AutoMigrate(
		&models.User{},
		&models.EmailOTP{},
		&models.Address{},
		&models.Category{},
		&models.Product{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate models: %w", err)
	}

	return database, nil
}

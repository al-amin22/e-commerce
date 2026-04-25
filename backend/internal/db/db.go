package db

import (
	"fmt"

	"ecommerce-backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(dbPath string) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := database.AutoMigrate(&models.User{}, &models.EmailOTP{}, &models.Address{}); err != nil {
		return nil, fmt.Errorf("failed to migrate models: %w", err)
	}

	return database, nil
}

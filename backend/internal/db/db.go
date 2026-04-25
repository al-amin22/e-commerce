package db

import (
	"fmt"

	"ecommerce-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(host, port, user, password, dbName, sslMode string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		host,
		port,
		user,
		password,
		dbName,
		sslMode,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := database.AutoMigrate(&models.User{}, &models.EmailOTP{}, &models.Address{}); err != nil {
		return nil, fmt.Errorf("failed to migrate models: %w", err)
	}

	return database, nil
}

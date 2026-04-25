package database

import (
	"fmt"

	"ecommerce-backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed connecting to postgres: %w", err)
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.Product{}); err != nil {
		return nil, fmt.Errorf("failed auto migrate: %w", err)
	}

	return db, nil
}

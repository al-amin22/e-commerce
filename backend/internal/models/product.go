package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string     `gorm:"size:120;not null;uniqueIndex" json:"name"`
	Slug      string     `gorm:"size:120;not null;uniqueIndex" json:"slug"`
	ParentID  *uuid.UUID `gorm:"type:char(36)" json:"parent_id"`
	Parent    *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children  []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type Product struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string     `gorm:"size:255;not null" json:"name"`
	Slug        string     `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Description string     `gorm:"type:text" json:"description"`
	Price       float64    `gorm:"not null" json:"price"`
	Stock       int        `gorm:"not null;default:0" json:"stock"`
	ImageURL    string     `gorm:"type:text" json:"image_url"`
	CategoryID  *uuid.UUID `gorm:"type:char(36)" json:"category_id"`
	Category    *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Weight      float64    `gorm:"default:0" json:"weight"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

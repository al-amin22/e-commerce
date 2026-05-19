package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleBuyer UserRole = "buyer"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name         string    `gorm:"size:120;not null" json:"name"`
	Email        string    `gorm:"size:200;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         UserRole  `gorm:"size:20;not null;default:buyer" json:"role"`
	IsVerified   bool      `gorm:"default:false" json:"is_verified"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	Addresses    []Address `json:"addresses,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.Role == "" {
		u.Role = RoleBuyer
	}
	return nil
}

type EmailOTP struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Email     string    `gorm:"index;not null" json:"email"`
	Code      string    `gorm:"size:10;not null" json:"code"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

func (o *EmailOTP) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

type Address struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"type:char(36);index;not null" json:"user_id"`
	Label        string    `gorm:"size:100;not null" json:"label"`
	Recipient    string    `gorm:"size:120;not null" json:"recipient"`
	Phone        string    `gorm:"size:30;not null" json:"phone"`
	AddressLine  string    `gorm:"type:text;not null" json:"address_line"`
	City         string    `gorm:"size:100;not null" json:"city"`
	Province     string    `gorm:"size:100;not null" json:"province"`
	PostalCode   string    `gorm:"size:20;not null" json:"postal_code"`
	CourierCode  string    `gorm:"size:30" json:"courier_code"`
	DestinationID string   `gorm:"size:50" json:"destination_id"`
	IsDefault    bool      `gorm:"default:false" json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

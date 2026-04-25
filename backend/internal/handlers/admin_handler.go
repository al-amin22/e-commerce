package handlers

import (
	"net/http"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ListUsers(c *gin.Context) {
	var users []models.User
	if err := h.DB.Select("id", "name", "email", "role", "is_verified", "created_at").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed getting users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *Handler) BootstrapAdmin(c *gin.Context) {
	var count int64
	h.DB.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "admin already exists"})
		return
	}

	passwordHash, _ := utils.HashPassword("admin123")
	admin := models.User{
		Name:         "Super Admin",
		Email:        "admin@local.dev",
		PasswordHash: passwordHash,
		Role:         models.RoleAdmin,
		IsVerified:   true,
	}

	if err := h.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed creating admin"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "admin created", "email": admin.Email, "password": "admin123"})
}

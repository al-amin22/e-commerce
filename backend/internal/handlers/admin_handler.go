package handlers

import (
	"net/http"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/response"
	"ecommerce-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ListUsers(c *gin.Context) {
	var users []models.User
	if err := h.DB.Select("id", "name", "email", "role", "is_verified", "created_at").Find(&users).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed getting users")
		return
	}
	response.Success(c, http.StatusOK, "success", gin.H{"users": users})
}

func (h *Handler) BootstrapAdmin(c *gin.Context) {
	var count int64
	h.DB.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count)
	if count > 0 {
		response.Error(c, http.StatusBadRequest, "admin already exists")
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
		response.Error(c, http.StatusInternalServerError, "failed creating admin")
		return
	}

	response.Success(c, http.StatusCreated, "admin created", gin.H{"email": admin.Email, "password": "admin123"})
}

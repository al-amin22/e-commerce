package handlers

import (
	"net/http"
	"strconv"

	"ecommerce-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createAddressRequest struct {
	Label         string `json:"label" binding:"required"`
	Recipient     string `json:"recipient" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
	AddressLine   string `json:"address_line" binding:"required"`
	City          string `json:"city" binding:"required"`
	Province      string `json:"province" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
	CourierCode   string `json:"courier_code"`
	DestinationID string `json:"destination_id"`
	IsDefault     bool   `json:"is_default"`
}

func (h *Handler) CreateAddress(c *gin.Context) {
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uuid.UUID)

	var req createAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if req.IsDefault {
		h.DB.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", false)
	}

	address := models.Address{
		UserID:        userID,
		Label:         req.Label,
		Recipient:     req.Recipient,
		Phone:         req.Phone,
		AddressLine:   req.AddressLine,
		City:          req.City,
		Province:      req.Province,
		PostalCode:    req.PostalCode,
		CourierCode:   req.CourierCode,
		DestinationID: req.DestinationID,
		IsDefault:     req.IsDefault,
	}

	if err := h.DB.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed creating address"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "address created", "address": address})
}

func (h *Handler) ListMyAddresses(c *gin.Context) {
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uuid.UUID)

	var addresses []models.Address
	if err := h.DB.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed getting addresses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"addresses": addresses})
}

func (h *Handler) GetShippingCost(c *gin.Context) {
	addressID := c.Query("address_id")
	weightStr := c.DefaultQuery("weight", "1000")
	courier := c.DefaultQuery("courier", "jne")

	weight, err := strconv.Atoi(weightStr)
	if err != nil || weight <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid weight"})
		return
	}

	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uuid.UUID)

	var address models.Address
	if err := h.DB.Where("id = ? AND user_id = ?", addressID, userID).First(&address).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "address not found"})
		return
	}

	if address.DestinationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "destination_id for this address is empty"})
		return
	}

	result, err := h.ShipSvc.GetCost(address.DestinationID, weight, courier)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"result":  result,
	})
}

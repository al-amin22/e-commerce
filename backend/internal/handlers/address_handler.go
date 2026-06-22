package handlers

import (
	"net/http"
	"strconv"

	"ecommerce-backend/internal/dto"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateAddress(c *gin.Context) {
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uuid.UUID)

	var req dto.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
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
		response.Error(c, http.StatusInternalServerError, "failed creating address")
		return
	}

	response.Success(c, http.StatusCreated, "address created", gin.H{"address": address})
}

func (h *Handler) ListMyAddresses(c *gin.Context) {
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uuid.UUID)

	var addresses []models.Address
	if err := h.DB.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed getting addresses")
		return
	}

	response.Success(c, http.StatusOK, "success", gin.H{"addresses": addresses})
}

func (h *Handler) GetShippingCost(c *gin.Context) {
	addressID := c.Query("address_id")
	weightStr := c.DefaultQuery("weight", "1000")
	courier := c.DefaultQuery("courier", "jne")

	weight, err := strconv.Atoi(weightStr)
	if err != nil || weight <= 0 {
		response.Error(c, http.StatusBadRequest, "invalid weight")
		return
	}

	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uuid.UUID)

	var address models.Address
	if err := h.DB.Where("id = ? AND user_id = ?", addressID, userID).First(&address).Error; err != nil {
		response.Error(c, http.StatusNotFound, "address not found")
		return
	}

	if address.DestinationID == "" {
		response.Error(c, http.StatusBadRequest, "destination_id for this address is empty")
		return
	}

	result, err := h.ShipSvc.GetCost(address.DestinationID, weight, courier)
	if err != nil {
		response.Error(c, http.StatusBadGateway, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "success", gin.H{
		"address": address,
		"result":  result,
	})
}

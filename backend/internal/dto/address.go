package dto

type CreateAddressRequest struct {
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

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ShippingService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

type ShippingCostResult struct {
	Courier string      `json:"courier"`
	Weight  int         `json:"weight"`
	Costs   interface{} `json:"costs"`
	Source  string      `json:"source"`
}

func NewShippingService(apiKey, baseURL string) *ShippingService {
	return &ShippingService{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *ShippingService) GetCost(destinationID string, weight int, courier string) (*ShippingCostResult, error) {
	if s.apiKey == "" {
		return &ShippingCostResult{
			Courier: courier,
			Weight:  weight,
			Source:  "mock",
			Costs: []map[string]interface{}{
				{"service": "REG", "description": "Regular", "cost": 18000, "etd": "2-3"},
				{"service": "YES", "description": "Yakin Esok Sampai", "cost": 32000, "etd": "1"},
			},
		}, nil
	}

	payload := map[string]string{
		"origin":      "501",
		"destination": destinationID,
		"weight":      strconv.Itoa(weight),
		"courier":     courier,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/cost", s.baseURL), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("key", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("shipping provider error: %s", resp.Status)
	}

	var parsed map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return &ShippingCostResult{
		Courier: courier,
		Weight:  weight,
		Source:  "rajaongkir",
		Costs:   parsed,
	}, nil
}

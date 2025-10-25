package infrastructure

import (
	"bytes"
	"encoding/json"
	"github.com/bekzxt/e-commerce/order-service/internal/domain"
	"net/http"
)

type InventoryClient struct {
	baseURL string
}

func NewInventoryClient(url string) *InventoryClient {
	return &InventoryClient{baseURL: url}
}

func (c *InventoryClient) CheckStock(items []domain.OrderItem) (bool, string, error) {
	reqBody := map[string]interface{}{
		"items": items,
	}

	body, _ := json.Marshal(reqBody)
	resp, err := http.Post(c.baseURL+"/products/check", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return false, "inventory service unavailable", err
	}
	defer resp.Body.Close()

	var result struct {
		Ok      bool    `json:"ok"`
		Message string  `json:"message"`
		Missing []int64 `json:"missing"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, "invalid inventory response", err
	}

	return result.Ok, result.Message, nil
}

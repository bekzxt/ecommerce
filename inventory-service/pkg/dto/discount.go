package dto

import (
	"github.com/bekzxt/e-commerce/inventory-service/internal/domain"
	"time"
)

type Discount struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	DiscountPercentage float64   `json:"discount-percentage"`
	ApplicableProducts []string  `json:"applicableProducts"`
	StartDate          time.Time `json:"start-date"`
	EndDate            time.Time `json:"end-date"`
	IsActive           bool      `json:"isActive"`
}
type DiscountedProduct struct {
	Product  domain.Product `json:"product"`
	Discount Discount       `json:"discount"`
}

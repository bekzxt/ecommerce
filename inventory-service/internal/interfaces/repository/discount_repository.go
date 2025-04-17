package repository

import "github.com/bekzxt/e-commerce/inventory-service/internal/domain"

type DiscountRepository interface {
	Create(product *domain.Discount) error
	Delete(id string) error
	GetDiscountedProducts() ([]domain.DiscountedProduct, error)
}

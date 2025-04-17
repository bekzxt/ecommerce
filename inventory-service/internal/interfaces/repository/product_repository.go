package repository

import "github.com/bekzxt/e-commerce/inventory-service/internal/domain"

type ProductRepository interface {
	Create(product *domain.Product) (*domain.Product, error)
	GetByID(id int64) (*domain.Product, error)
	Update(product *domain.Product) (*domain.Product, error)
	Delete(id int64) error
	List() ([]*domain.Product, error)
}

package usecase

import "github.com/bekzxt/e-commerce/inventory-service/internal/domain"

type Inventory interface {
	CreateUC(p *domain.Product) (*domain.Product, error)
	GetByID(id int64) (*domain.Product, error)
	Delete(id int64) error
	List() ([]*domain.Product, error)
	ReserveItems(items []domain.OrderItemInv) (bool, error)
}

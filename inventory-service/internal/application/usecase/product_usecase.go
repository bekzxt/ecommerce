package usecase

import (
	"github.com/bekzxt/e-commerce/inventory-service/internal/domain"
	"github.com/bekzxt/e-commerce/inventory-service/internal/interfaces/repository"
)

type ProductUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: repo}
}

func (uc *ProductUseCase) CreateUC(p *domain.Product) (*domain.Product, error) {
	return uc.repo.Create(p)
}

func (uc *ProductUseCase) GetByID(id int64) (*domain.Product, error) {
	return uc.repo.GetByID(id)
}

func (uc *ProductUseCase) Update(p *domain.Product) (*domain.Product, error) {
	return uc.repo.Update(p)
}

func (uc *ProductUseCase) Delete(id int64) error {
	return uc.repo.Delete(id)
}

func (uc *ProductUseCase) List() ([]*domain.Product, error) {
	return uc.repo.List()
}

package usecase

import (
	"github.com/bekzxt/e-commerce/inventory-service/internal/domain"
	"github.com/bekzxt/e-commerce/inventory-service/internal/interfaces/repository"
)

type DiscountUseCase struct {
	repo repository.DiscountRepository
}

func NewDiscountUseCase(repo repository.DiscountRepository) *DiscountUseCase {
	return &DiscountUseCase{repo: repo}
}

func (uc *DiscountUseCase) Create(p *domain.Discount) (*domain.Discount, error) {
	return nil, uc.repo.Create(p)
}

func (uc *DiscountUseCase) Delete(id string) error {
	return uc.repo.Delete(id)
}

func (uc *DiscountUseCase) List() ([]domain.DiscountedProduct, error) {
	return uc.repo.GetDiscountedProducts()
}

package usecase

import (
	"github.com/bekzxt/e-commerce/order-service/internal/domain"
	"github.com/bekzxt/e-commerce/order-service/internal/interfaces/repository"
)

type ReviewUseCase struct {
	reviewRepo repository.ReviewRepository
}

func NewReviewUseCase(repo repository.ReviewRepository) *ReviewUseCase {
	return &ReviewUseCase{reviewRepo: repo}
}

func (uc *ReviewUseCase) Create(p *domain.Review) (*domain.Review, error) {
	return uc.reviewRepo.Create(p)
}
func (uc *ReviewUseCase) Update(p *domain.Review) (*domain.Review, error) {
	return uc.reviewRepo.Update(p)
}

package repository

import "github.com/bekzxt/e-commerce/order-service/internal/domain"

type ReviewRepository interface {
	Create(review *domain.Review) (*domain.Review, error)
	Update(review *domain.Review) (*domain.Review, error)
}

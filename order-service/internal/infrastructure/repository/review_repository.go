package repository

import (
	"database/sql"
	"github.com/bekzxt/e-commerce/order-service/internal/domain"
	"github.com/bekzxt/e-commerce/order-service/internal/interfaces/repository"
)

type ReviewRepositoryImpl struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) repository.ReviewRepository {
	return &ReviewRepositoryImpl{db: db}
}

func (r *ReviewRepositoryImpl) Create(review *domain.Review) (*domain.Review, error) {
	err := r.db.QueryRow(`
		INSERT INTO reviews (product_id, user_id, comment, rating)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, review.ProductID, review.UserID, review.Comment, review.Rating).Scan(&review.ID)

	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *ReviewRepositoryImpl) Update(review *domain.Review) (*domain.Review, error) {
	_, err := r.db.Exec(`
		UPDATE reviews
		SET comment=$1, rating=$2
		WHERE id=$3
	`, review.Comment, review.Rating, review.ID)
	if err != nil {
		return nil, err
	}
	var updated domain.Review
	err = r.db.QueryRow(`
		SELECT id, product_id, user_id, comment, rating
		FROM reviews
		WHERE id=$1
	`, review.ID).Scan(&updated.ID, &updated.ProductID, &updated.UserID, &updated.Comment, &updated.Rating)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

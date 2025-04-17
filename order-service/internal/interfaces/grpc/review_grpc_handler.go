package grpc

import (
	"github.com/bekzxt/e-commerce/order-service/internal/application/usecase"
	"github.com/bekzxt/e-commerce/order-service/internal/domain"
	pbb "github.com/bekzxt/e-commerce/order-service/proto_review"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReviewHandler struct {
	pbb.UnimplementedReviewServiceServer
	usecase *usecase.ReviewUseCase
}

func NewReviewHandler(uc *usecase.ReviewUseCase) *ReviewHandler {
	return &ReviewHandler{usecase: uc}
}

func (h *ReviewHandler) CreateReview(ctx context.Context, req *pbb.CreateReviewRequest) (*pbb.ReviewResponse, error) {
	p := &domain.Review{
		ProductID: req.Review.ProductId,
		UserID:    req.Review.UserId,
		Rating:    req.Review.Rating,
		Comment:   req.Review.Comment,
	}
	created, err := h.usecase.Create(p)
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, status.Error(codes.Internal, "failed to create product: nil response")
	}
	return &pbb.ReviewResponse{Review: toProto(created)}, nil
}

func (h *ReviewHandler) UpdateReview(ctx context.Context, req *pbb.UpdateReviewRequest) (*pbb.ReviewResponse, error) {
	p := &domain.Review{
		ID:        req.Review.Id,
		ProductID: req.Review.ProductId,
		UserID:    req.Review.UserId,
		Rating:    req.Review.Rating,
		Comment:   req.Review.Comment,
	}
	updated, err := h.usecase.Update(p)
	if err != nil {
		return nil, err
	}
	return &pbb.ReviewResponse{Review: toProto(updated)}, nil
}

func toProto(p *domain.Review) *pbb.Review {
	return &pbb.Review{
		Id:        p.ID,
		Rating:    p.Rating,
		UserId:    p.UserID,
		Comment:   p.Comment,
		ProductId: p.ProductID,
	}
}

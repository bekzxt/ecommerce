package grpc

import (
	"context"
	"github.com/bekzxt/e-commerce/inventory-service/internal/application/usecase"
	"github.com/bekzxt/e-commerce/inventory-service/internal/domain"
	pb "github.com/bekzxt/e-commerce/inventory-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	pb.UnimplementedInventoryServiceServer
	uc *usecase.ProductUseCase
}

func NewProductHandler(svc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc: svc}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	p := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Stock:       req.Stock,
		Price:       req.Price,
		CategoryID:  req.CategoryId,
	}
	created, err := h.uc.CreateUC(p)
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, status.Error(codes.Internal, "failed to create product: nil response")
	}
	return &pb.ProductResponse{Product: toProto(created)}, nil
}
func (h *ProductHandler) GetProductByID(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	product, err := h.uc.GetByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ProductResponse{Product: toProto(product)}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	p := &domain.Product{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Stock:       req.Stock,
		Price:       req.Price,
		CategoryID:  req.CategoryId,
	}
	updated, err := h.uc.Update(p)
	if err != nil {
		return nil, err
	}
	return &pb.ProductResponse{Product: toProto(updated)}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.Empty, error) {
	err := h.uc.Delete(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := h.uc.List()
	if err != nil {
		return nil, err
	}
	var list []*pb.Product
	for _, p := range products {
		list = append(list, toProto(p))
	}
	return &pb.ListProductsResponse{Products: list}, nil
}

func toProto(p *domain.Product) *pb.Product {
	return &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Stock:       p.Stock,
		Price:       p.Price,
		CategoryId:  p.CategoryID,
	}
}

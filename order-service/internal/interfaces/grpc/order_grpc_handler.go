package grpc

import (
	"context"

	"github.com/bekzxt/e-commerce/order-service/internal/application/usecase"
	"github.com/bekzxt/e-commerce/order-service/internal/domain"
	"github.com/bekzxt/e-commerce/order-service/internal/interfaces/dto"
	pb "github.com/bekzxt/e-commerce/order-service/proto"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	usecase *usecase.OrderUseCase
}

func NewOrderHandler(uc *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{usecase: uc}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	dtoReq := dto.CreateOrderRequest{
		UserID: req.UserId,
	}

	for _, item := range req.Items {
		dtoReq.Items = append(dtoReq.Items, dto.OrderItemRequest{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
			Price:     item.Price,
		})
	}

	order, err := h.usecase.CreateOrder(dtoReq)
	if err != nil {
		return nil, err
	}

	return toProtoOrder(order), nil
}

func (h *OrderHandler) GetOrderByID(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	order, err := h.usecase.GetOrderByID(req.OrderId)
	if err != nil {
		return nil, err
	}
	return toProtoOrder(order), nil
}

func (h *OrderHandler) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error) {
	err := h.usecase.UpdateOrderStatus(req.OrderId, domain.OrderStatus(req.Status))
	if err != nil {
		return nil, err
	}

	order, err := h.usecase.GetOrderByID(req.OrderId)
	if err != nil {
		return nil, err
	}
	return toProtoOrder(order), nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := h.usecase.ListOrdersByUser(req.UserId)
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.OrderResponse
	for _, order := range orders {
		pbOrders = append(pbOrders, toProtoOrder(order))
	}

	return &pb.ListOrdersResponse{
		Orders: pbOrders,
	}, nil
}

func toProtoOrder(o *domain.Order) *pb.OrderResponse {
	var pbItems []*pb.OrderItem
	for _, item := range o.Items {
		pbItems = append(pbItems, &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		})
	}

	return &pb.OrderResponse{
		Id:         o.ID,
		UserId:     o.UserID,
		Status:     string(o.Status),
		TotalPrice: o.TotalPrice,
		Items:      pbItems,
	}
}

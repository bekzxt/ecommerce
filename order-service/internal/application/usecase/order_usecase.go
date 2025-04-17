package usecase

import (
	"github.com/bekzxt/e-commerce/order-service/internal/domain"
	"github.com/bekzxt/e-commerce/order-service/internal/interfaces/dto"
	"github.com/bekzxt/e-commerce/order-service/internal/interfaces/repository"
	"github.com/google/uuid"
)

type OrderUseCase struct {
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
}

func NewOrderUseCase(repo repository.OrderRepository, i repository.OrderItemRepository) *OrderUseCase {
	return &OrderUseCase{orderRepo: repo, orderItemRepo: i}
}

func (uc *OrderUseCase) CreateOrder(req dto.CreateOrderRequest) (*domain.Order, error) {
	orderID := uuid.New().String()
	var items []domain.OrderItem
	var total float64
	for _, item := range req.Items {
		items = append(items, domain.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
		total += float64(item.Quantity) * item.Price
	}
	order := &domain.Order{
		ID:         orderID,
		UserID:     req.UserID,
		TotalPrice: total,
		Items:      items,
		Status:     domain.Pending,
	}
	if err := uc.orderRepo.CreateOrder(order); err != nil {
		return nil, err
	}

	for _, item := range items {
		if err := uc.orderItemRepo.CreateOrderItem(&item); err != nil {
			return nil, err
		}
	}
	return order, nil
}

func (uc *OrderUseCase) GetOrderByID(orderID string) (*domain.Order, error) {
	return uc.orderRepo.GetOrderByID(orderID)
}

func (uc *OrderUseCase) UpdateOrderStatus(orderID string, newStatus domain.OrderStatus) error {
	if newStatus != domain.Pending && newStatus != domain.Completed && newStatus != domain.Cancelled {
		return domain.ErrInvalidStatus
	}
	return uc.orderRepo.UpdateOrderStatus(orderID, newStatus)
}

func (uc *OrderUseCase) ListOrdersByUser(userID string) ([]*domain.Order, error) {
	return uc.orderRepo.ListOrdersByUserID(userID)
}

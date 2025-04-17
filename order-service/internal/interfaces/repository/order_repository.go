package repository

import "github.com/bekzxt/e-commerce/order-service/internal/domain"

type OrderRepository interface {
	CreateOrder(order *domain.Order) error
	GetOrderByID(id string) (*domain.Order, error)
	UpdateOrderStatus(id string, status domain.OrderStatus) error
	ListOrdersByUserID(userID string) ([]*domain.Order, error)
}

type OrderItemRepository interface {
	CreateOrderItem(item *domain.OrderItem) error
	GetItemsByOrderID(orderID string) ([]domain.OrderItem, error)
	DeleteItemsByOrderID(orderID string) error
}

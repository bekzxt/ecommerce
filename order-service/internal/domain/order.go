package domain

import "errors"

type Order struct {
	ID         string
	UserID     string
	Status     OrderStatus
	TotalPrice float64
	Items      []OrderItem
}
type OrderItem struct {
	OrderID   string  `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Completed OrderStatus = "completed"
	Cancelled OrderStatus = "cancelled"
)

type OrderCreatedEvent struct {
	OrderID string      `json:"order_id"`
	UserID  string      `json:"user_id"`
	Total   float64     `json:"total"`
	Status  string      `json:"status"`
	Items   []OrderItem `json:"items"`
}

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidStatus = errors.New("invalid status")
)

type OrderReservedEvent struct {
	OrderID string `json:"order_id"`
}

type OrderRejectedEvent struct {
	OrderID string `json:"order_id"`
	Reason  string `json:"reason,omitempty"`
}

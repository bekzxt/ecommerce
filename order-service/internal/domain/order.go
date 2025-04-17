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
	OrderID   string
	ProductID string
	Quantity  int
	Price     float64
}

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Completed OrderStatus = "completed"
	Cancelled OrderStatus = "cancelled"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidStatus = errors.New("invalid status")
)

package dto

type CreateOrderRequest struct {
	UserID     string             `json:"user_id" binding:"required"`
	TotalPrice float64            `json:"total_price"`
	Items      []OrderItemRequest `json:"items" binding:"required"`
}

type OrderItemRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	Price     float64 `json:"price" binding:"required,gt=0"`
}

type OrderResponse struct {
	ID         string              `json:"id"`
	UserID     string              `json:"user_id"`
	Status     string              `json:"status"`
	TotalPrice float64             `json:"total_price"`
	Items      []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	ProductID string  `json:"productID"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

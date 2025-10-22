package domain

type OrderItemInv struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       float64
	Stock       int32
	CategoryID  int32
}

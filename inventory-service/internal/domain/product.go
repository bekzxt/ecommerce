package domain

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       float64
	Stock       int32
	CategoryID  int32
}

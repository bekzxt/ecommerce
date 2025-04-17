package domain

import "time"

type Discount struct {
	ID                 string
	Name               string
	Description        string
	DiscountPercentage float64
	ApplicableProducts []string
	StartDate          time.Time
	EndDate            time.Time
	IsActive           bool
}

type DiscountedProduct struct {
	Product  Product
	Discount Discount
}

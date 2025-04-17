package domain

import "time"

type Review struct {
	ID        uint64
	ProductID string
	UserID    string
	Rating    float64
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

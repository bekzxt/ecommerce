package repository

import (
	"database/sql"
	"github.com/bekzxt/e-commerce/order-service/internal/domain"
)

type OrderItemRepositoryImpl struct {
	db *sql.DB
}

func NewOrderItemRepository(db *sql.DB) *OrderItemRepositoryImpl {
	return &OrderItemRepositoryImpl{db: db}
}

func (r *OrderItemRepositoryImpl) CreateOrderItem(item *domain.OrderItem) error {
	query := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, item.OrderID, item.ProductID, item.Quantity, item.Price)
	return err
}

func (r *OrderItemRepositoryImpl) GetItemsByOrderID(orderID string) ([]domain.OrderItem, error) {
	query := `SELECT order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`
	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *OrderItemRepositoryImpl) DeleteItemsByOrderID(orderID string) error {
	query := `DELETE FROM order_items WHERE order_id = $1`
	_, err := r.db.Exec(query, orderID)
	return err
}

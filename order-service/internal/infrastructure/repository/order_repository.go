package repository

import (
	"database/sql"
	"errors"
	"github.com/bekzxt/e-commerce/order-service/internal/domain"
	"github.com/bekzxt/e-commerce/order-service/internal/interfaces/repository"
)

type OrderRepositoryImpl struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

func (r *OrderRepositoryImpl) CreateOrder(order *domain.Order) error {
	query := `INSERT INTO orders (id, user_id, total_price, status) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, order.ID, order.UserID, order.TotalPrice, order.Status)
	return err
}
func (r *OrderRepositoryImpl) GetOrderByID(id string) (*domain.Order, error) {
	query := `SELECT id, user_id, total_price, status FROM orders WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var order domain.Order
	if err := row.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	orderItemsQuery := `SELECT product_id, quantity, price FROM order_items WHERE order_id = $1`
	rows, err := r.db.Query(orderItemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}

func (r *OrderRepositoryImpl) UpdateOrderStatus(id string, status domain.OrderStatus) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}
func (r *OrderRepositoryImpl) ListOrdersByUserID(userID string) ([]*domain.Order, error) {
	query := `SELECT id, user_id, total_price, status FROM orders WHERE user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status); err != nil {
			return nil, err
		}

		orderItemsQuery := `SELECT product_id, quantity, price FROM order_items WHERE order_id = $1`
		rows, err1 := r.db.Query(orderItemsQuery, order.ID)
		if err1 != nil {
			return nil, err1
		}
		defer rows.Close()

		for rows.Next() {
			var item domain.OrderItem
			if err := rows.Scan(&item.ProductID, &item.Quantity, &item.Price); err != nil {
				return nil, err
			}
			order.Items = append(order.Items, item)
		}

		orders = append(orders, &order)

	}

	return orders, nil
}

package repository

import (
	"database/sql"
	"github.com/bekzxt/e-commerce/inventory-service/internal/domain"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db}
}

func (r *ProductRepo) Create(p *domain.Product) (*domain.Product, error) {
	err := r.db.QueryRow(`
		INSERT INTO products (name, description, price, stock, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, p.Name, p.Description, p.Price, p.Stock, p.CategoryID).Scan(&p.ID)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *ProductRepo) GetByID(id int64) (*domain.Product, error) {
	row := r.db.QueryRow(`SELECT id, name, description, price, stock, category_id FROM products WHERE id=$1`, id)

	var p domain.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CategoryID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepo) Update(p *domain.Product) (*domain.Product, error) {
	_, err := r.db.Exec(`UPDATE products SET name=$1, description=$2, price=$3, stock=$4, category_id=$5 WHERE id=$6`,
		p.Name, p.Description, p.Price, p.Stock, p.CategoryID, p.ID)
	if err != nil {
		return nil, err
	}

	var updated domain.Product
	err = r.db.QueryRow(`SELECT id, name, description, price, stock, category_id FROM products WHERE id=$1`, p.ID).
		Scan(&updated.ID, &updated.Name, &updated.Description, &updated.Price, &updated.Stock, &updated.CategoryID)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *ProductRepo) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM products WHERE id=$1`, id)
	return err
}

func (r *ProductRepo) List() ([]*domain.Product, error) {
	rows, err := r.db.Query(`SELECT id, name, description, price, stock, category_id FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

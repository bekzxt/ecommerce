package repository

import (
	"database/sql"
	"github.com/bekzxt/e-commerce/inventory-service/internal/domain"
	"github.com/lib/pq"
)

type DiscountRepo struct {
	db *sql.DB
}

type DiscountedProduct struct {
	Product  domain.Product
	Discount domain.Discount
}

func NewDiscountRepo(db *sql.DB) *DiscountRepo {
	return &DiscountRepo{db}
}

func (r *DiscountRepo) Create(p *domain.Discount) error {
	err := r.db.QueryRow(`
    INSERT INTO discounts (name, description, discount_percentage, applicable_products, start_date, end_date, is_active)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id`,
		p.Name, p.Description, p.DiscountPercentage, pq.Array(p.ApplicableProducts), p.StartDate, p.EndDate, p.IsActive).Scan(&p.ID)
	return err
}

func (r *DiscountRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM discounts WHERE id=$1`, id)
	return err
}

func (r *DiscountRepo) GetDiscountedProducts() ([]domain.DiscountedProduct, error) {
	rows, err := r.db.Query(`
		SELECT 
			p.id, p.name, p.description, p.price, p.stock, p.category_id,
			d.id, d.name, d.description, d.discount_percentage, d.applicable_products, d.start_date, d.end_date, d.is_active
		FROM discounts d
		JOIN unnest(d.applicable_products::integer[]) AS product_id ON true
		JOIN products p ON p.id = product_id
		WHERE d.is_active = true
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.DiscountedProduct
	for rows.Next() {
		var prod domain.Product
		var disc domain.Discount
		err := rows.Scan(
			&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Stock, &prod.CategoryID,
			&disc.ID, &disc.Name, &disc.Description, &disc.DiscountPercentage, pq.Array(&disc.ApplicableProducts), &disc.StartDate, &disc.EndDate, &disc.IsActive,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, domain.DiscountedProduct{
			Product:  prod,
			Discount: disc,
		})
	}

	return result, nil
}

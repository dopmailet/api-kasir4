package repositories

import (
	"database/sql"
	"kasir-api/models"
)

// DiscountRepository handles database operations for discounts
type DiscountRepository struct {
	db *sql.DB
}

// NewDiscountRepository creates a new DiscountRepository
func NewDiscountRepository(db *sql.DB) *DiscountRepository {
	return &DiscountRepository{db: db}
}

// Create inserts a new discount into the database
func (r *DiscountRepository) Create(d *models.Discount) error {
	query := `
		INSERT INTO discounts (name, type, value, min_order_amount, start_date, end_date, is_active, product_id, category_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	return r.db.QueryRow(query, d.Name, d.Type, d.Value, d.MinOrderAmount, d.StartDate, d.EndDate, d.IsActive, d.ProductID, d.CategoryID).Scan(&d.ID)
}

// GetAll returns all discounts (for admin management)
func (r *DiscountRepository) GetAll() ([]models.Discount, error) {
	query := `
		SELECT id, name, type, value, min_order_amount, start_date, end_date, is_active, product_id, category_id 
		FROM discounts 
		ORDER BY start_date DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var discounts []models.Discount
	for rows.Next() {
		var d models.Discount
		if err := rows.Scan(&d.ID, &d.Name, &d.Type, &d.Value, &d.MinOrderAmount, &d.StartDate, &d.EndDate, &d.IsActive, &d.ProductID, &d.CategoryID); err != nil {
			return nil, err
		}
		discounts = append(discounts, d)
	}
	return discounts, nil
}

// GetActive returns only active and valid GLOBAL discounts (for cashier selection)
// Product and Category discounts are applied automatically, not selected manually.
// So this should return ONLY Global discounts (product_id IS NULL AND category_id IS NULL)
func (r *DiscountRepository) GetActive() ([]models.Discount, error) {
	query := `
		SELECT id, name, type, value, min_order_amount, start_date, end_date, is_active, product_id, category_id
		FROM discounts 
		WHERE is_active = TRUE 
		AND product_id IS NULL
		AND category_id IS NULL
		AND NOW() BETWEEN start_date AND end_date
		ORDER BY min_order_amount ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var discounts []models.Discount
	for rows.Next() {
		var d models.Discount
		if err := rows.Scan(&d.ID, &d.Name, &d.Type, &d.Value, &d.MinOrderAmount, &d.StartDate, &d.EndDate, &d.IsActive, &d.ProductID, &d.CategoryID); err != nil {
			return nil, err
		}
		discounts = append(discounts, d)
	}
	return discounts, nil
}

// GetByID returns a discount by ID
func (r *DiscountRepository) GetByID(id int) (*models.Discount, error) {
	query := `SELECT id, name, type, value, min_order_amount, start_date, end_date, is_active, product_id, category_id FROM discounts WHERE id = $1`
	var d models.Discount
	err := r.db.QueryRow(query, id).Scan(&d.ID, &d.Name, &d.Type, &d.Value, &d.MinOrderAmount, &d.StartDate, &d.EndDate, &d.IsActive, &d.ProductID, &d.CategoryID)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Update updates an existing discount
func (r *DiscountRepository) Update(id int, d *models.Discount) error {
	query := `
		UPDATE discounts 
		SET name=$1, type=$2, value=$3, min_order_amount=$4, start_date=$5, end_date=$6, is_active=$7, product_id=$8, category_id=$9
		WHERE id=$10
	`
	_, err := r.db.Exec(query, d.Name, d.Type, d.Value, d.MinOrderAmount, d.StartDate, d.EndDate, d.IsActive, d.ProductID, d.CategoryID, id)
	return err
}

// Delete removes a discount
func (r *DiscountRepository) Delete(id int) error {
	query := `DELETE FROM discounts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

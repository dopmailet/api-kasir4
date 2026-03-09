package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"strings"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(c *models.Customer) error {
	query := `
		INSERT INTO customers (name, phone, card_number, address, notes, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, loyalty_points, total_spent, total_transactions, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		c.Name, c.Phone, c.CardNumber, c.Address, c.Notes, c.IsActive,
	).Scan(&c.ID, &c.LoyaltyPoints, &c.TotalSpent, &c.TotalTransactions, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CustomerRepository) GetByID(id int) (*models.Customer, error) {
	query := `
		SELECT 
			c.id, c.name, c.phone, c.card_number, c.address, c.notes, c.loyalty_points, 
			c.total_spent, c.total_transactions, c.last_transaction_at, c.is_active, c.created_at, c.updated_at,
			COALESCE(
				(SELECT SUM(td.quantity * (td.price - COALESCE(td.harga_beli, 0)) - td.discount_amount)
				 FROM transactions t
				 JOIN transaction_details td ON td.transaction_id = t.id
				 WHERE t.customer_id = c.id),
				0
			) AS total_profit
		FROM customers c WHERE c.id = $1
	`
	var c models.Customer
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.Name, &c.Phone, &c.CardNumber, &c.Address, &c.Notes, &c.LoyaltyPoints,
		&c.TotalSpent, &c.TotalTransactions, &c.LastTransactionAt, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
		&c.TotalProfit,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("customer id %d tidak ditemukan", id)
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepository) GetAll(search string, status string, page, limit int, sortBy, sortOrder string) ([]models.Customer, int, error) {
	// Base query
	whereClauses := []string{}
	args := []interface{}{}
	argId := 1

	if search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(c.name ILIKE $%d OR c.phone ILIKE $%d OR c.card_number ILIKE $%d)", argId, argId, argId))
		args = append(args, "%"+search+"%")
		argId++
	}

	if status == "active" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.is_active = $%d", argId))
		args = append(args, true)
		argId++
	} else if status == "inactive" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.is_active = $%d", argId))
		args = append(args, false)
		argId++
	}

	whereStmt := ""
	if len(whereClauses) > 0 {
		whereStmt = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM customers c " + whereStmt
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Sort checking
	validSortFields := map[string]string{
		"id":                  "c.id",
		"name":                "c.name",
		"loyalty_points":      "c.loyalty_points",
		"total_spent":         "c.total_spent",
		"last_transaction_at": "c.last_transaction_at",
		"created_at":          "c.created_at",
	}

	dbSortField, valid := validSortFields[sortBy]
	if !valid {
		dbSortField = "c.created_at"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Pagination
	offset := (page - 1) * limit
	query := fmt.Sprintf(`
		SELECT 
			c.id, c.name, c.phone, c.card_number, c.address, c.notes, c.loyalty_points, 
			c.total_spent, c.total_transactions, c.last_transaction_at, c.is_active, c.created_at, c.updated_at,
			COALESCE(
				(SELECT SUM(td.quantity * (td.price - COALESCE(td.harga_beli, 0)) - td.discount_amount)
				 FROM transactions t
				 JOIN transaction_details td ON td.transaction_id = t.id
				 WHERE t.customer_id = c.id),
				0
			) AS total_profit
		FROM customers c
		%s
		ORDER BY %s %s, c.id DESC
		LIMIT $%d OFFSET $%d
	`, whereStmt, dbSortField, sortOrder, argId, argId+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		err := rows.Scan(
			&c.ID, &c.Name, &c.Phone, &c.CardNumber, &c.Address, &c.Notes, &c.LoyaltyPoints,
			&c.TotalSpent, &c.TotalTransactions, &c.LastTransactionAt, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
			&c.TotalProfit,
		)
		if err != nil {
			return nil, 0, err
		}
		customers = append(customers, c)
	}

	return customers, total, nil
}

func (r *CustomerRepository) Update(c *models.Customer) error {
	query := `
		UPDATE customers 
		SET name = $1, phone = $2, card_number = $3, address = $4, notes = $5, is_active = $6
		WHERE id = $7
		RETURNING updated_at
	`
	err := r.db.QueryRow(
		query, c.Name, c.Phone, c.CardNumber, c.Address, c.Notes, c.IsActive, c.ID,
	).Scan(&c.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("customer id %d tidak ditemukan", c.ID)
	}
	return err
}

func (r *CustomerRepository) GenerateCustomerCode() (string, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM customers").Scan(&count)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("CUST-%06d", count+1), nil
}

// GetTransactions retrieves the order history for a specific customer
func (r *CustomerRepository) GetTransactions(customerID int) ([]models.TransactionWithItems, error) {
	query := `
		SELECT t.id, t.created_at, t.total_amount, u.username
		FROM transactions t
		LEFT JOIN users u ON t.created_by = u.id
		WHERE t.customer_id = $1
		ORDER BY t.created_at DESC
		LIMIT 50
	`
	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []models.TransactionWithItems
	for rows.Next() {
		var t models.TransactionWithItems
		var username sql.NullString
		err := rows.Scan(&t.ID, &t.CreatedAt, &t.TotalAmount, &username)
		if err != nil {
			return nil, err
		}
		if username.Valid {
			t.Username = username.String
		}
		txs = append(txs, t)
	}
	return txs, nil
}

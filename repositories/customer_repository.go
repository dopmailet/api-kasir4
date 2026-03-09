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
		INSERT INTO customers (customer_code, name, phone, address, notes, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, loyalty_points, total_spent, total_transactions, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		c.CustomerCode, c.Name, c.Phone, c.Address, c.Notes, c.IsActive,
	).Scan(&c.ID, &c.LoyaltyPoints, &c.TotalSpent, &c.TotalTransactions, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CustomerRepository) GetByID(id int) (*models.Customer, error) {
	query := `
		SELECT id, customer_code, name, phone, address, notes, loyalty_points, 
			     total_spent, total_transactions, last_transaction_at, is_active, created_at, updated_at
		FROM customers WHERE id = $1
	`
	var c models.Customer
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.CustomerCode, &c.Name, &c.Phone, &c.Address, &c.Notes, &c.LoyaltyPoints,
		&c.TotalSpent, &c.TotalTransactions, &c.LastTransactionAt, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
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
		whereClauses = append(whereClauses, fmt.Sprintf("(name ILIKE $%d OR phone ILIKE $%d OR customer_code ILIKE $%d)", argId, argId, argId))
		args = append(args, "%"+search+"%")
		argId++
	}

	if status == "active" {
		whereClauses = append(whereClauses, fmt.Sprintf("is_active = $%d", argId))
		args = append(args, true)
		argId++
	} else if status == "inactive" {
		whereClauses = append(whereClauses, fmt.Sprintf("is_active = $%d", argId))
		args = append(args, false)
		argId++
	}

	whereStmt := ""
	if len(whereClauses) > 0 {
		whereStmt = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM customers " + whereStmt
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Sort checking
	validSortFields := map[string]string{
		"id":                  "id",
		"name":                "name",
		"loyalty_points":      "loyalty_points",
		"total_spent":         "total_spent",
		"last_transaction_at": "last_transaction_at",
		"created_at":          "created_at",
	}

	dbSortField, valid := validSortFields[sortBy]
	if !valid {
		dbSortField = "created_at"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Pagination
	offset := (page - 1) * limit
	query := fmt.Sprintf(`
		SELECT id, customer_code, name, phone, address, notes, loyalty_points, 
			     total_spent, total_transactions, last_transaction_at, is_active, created_at, updated_at
		FROM customers 
		%s
		ORDER BY %s %s, id DESC
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
			&c.ID, &c.CustomerCode, &c.Name, &c.Phone, &c.Address, &c.Notes, &c.LoyaltyPoints,
			&c.TotalSpent, &c.TotalTransactions, &c.LastTransactionAt, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
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
		SET name = $1, phone = $2, address = $3, notes = $4, is_active = $5
		WHERE id = $6
		RETURNING updated_at
	`
	err := r.db.QueryRow(
		query, c.Name, c.Phone, c.Address, c.Notes, c.IsActive, c.ID,
	).Scan(&c.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("customer id %d tidak ditemukan", c.ID)
	}
	return err
}

func (r *CustomerRepository) GenerateCustomerCode() (string, error) {
	// Simple code generator: CUST-00000X
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM customers").Scan(&count)
	if err != nil {
		return "", err
	}
	// Prefix format e.g: CUST-000001
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

package repositories

import (
	"database/sql"
	"kasir-api/models"
	"strconv"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// GetAll mengambil semua data pengeluaran dengan optional filter bulan/tahun
func (r *ExpenseRepository) GetAll(year string, month string) ([]models.Expense, error) {
	query := `
		SELECT 
			e.id, 
			e.category, 
			e.description, 
			e.amount, 
			TO_CHAR(e.expense_date, 'YYYY-MM-DD') AS expense_date, 
			e.is_recurring, 
			e.recurring_period, 
			e.notes, 
			e.created_by, 
			u.username as creator_name,
			e.created_at, 
			e.updated_at
		FROM expenses e
		LEFT JOIN users u ON e.created_by = u.id
		WHERE 1=1
	`
	var args []interface{}
	argCount := 1

	if year != "" {
		query += ` AND EXTRACT(YEAR FROM e.expense_date) = $` + strconv.Itoa(argCount)
		args = append(args, year)
		argCount++
	}

	if month != "" {
		query += ` AND EXTRACT(MONTH FROM e.expense_date) = $` + strconv.Itoa(argCount)
		args = append(args, month)
		argCount++
	}

	query += ` ORDER BY e.expense_date DESC, e.created_at DESC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		var recPeriod, notes sql.NullString
		var creatorName sql.NullString

		if err := rows.Scan(
			&e.ID, &e.Category, &e.Description, &e.Amount, &e.ExpenseDate,
			&e.IsRecurring, &recPeriod, &notes, &e.CreatedBy, &creatorName,
			&e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if recPeriod.Valid {
			e.RecurringPeriod = &recPeriod.String
		}
		if notes.Valid {
			e.Notes = &notes.String
		}
		if creatorName.Valid {
			e.CreatorName = creatorName.String
		}

		expenses = append(expenses, e)
	}

	if expenses == nil {
		expenses = []models.Expense{}
	}

	return expenses, nil
}

// GetByID mengambil satu data expense
func (r *ExpenseRepository) GetByID(id int) (*models.Expense, error) {
	query := `
		SELECT 
			e.id, e.category, e.description, e.amount, 
			TO_CHAR(e.expense_date, 'YYYY-MM-DD') AS expense_date, 
			e.is_recurring, e.recurring_period, e.notes, e.created_by,
			e.created_at, e.updated_at
		FROM expenses e
		WHERE e.id = $1
	`
	var e models.Expense
	var recPeriod, notes sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&e.ID, &e.Category, &e.Description, &e.Amount, &e.ExpenseDate,
		&e.IsRecurring, &recPeriod, &notes, &e.CreatedBy,
		&e.CreatedAt, &e.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if recPeriod.Valid {
		e.RecurringPeriod = &recPeriod.String
	}
	if notes.Valid {
		e.Notes = &notes.String
	}

	return &e, nil
}

// Create menambahkan data expense baru
func (r *ExpenseRepository) Create(e *models.Expense) (*models.Expense, error) {
	query := `
		INSERT INTO expenses 
			(category, description, amount, expense_date, is_recurring, recurring_period, notes, created_by)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING 
			id, TO_CHAR(expense_date, 'YYYY-MM-DD'), created_at, updated_at
	`
	err := r.db.QueryRow(
		query,
		e.Category, e.Description, e.Amount, e.ExpenseDate,
		e.IsRecurring, e.RecurringPeriod, e.Notes, e.CreatedBy,
	).Scan(&e.ID, &e.ExpenseDate, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return e, nil
}

// Update memodifikasi data expense yang ada
func (r *ExpenseRepository) Update(id int, e *models.Expense) (*models.Expense, error) {
	query := `
		UPDATE expenses
		SET 
			category = $1, 
			description = $2, 
			amount = $3, 
			expense_date = $4, 
			is_recurring = $5, 
			recurring_period = $6, 
			notes = $7,
			updated_at = NOW()
		WHERE id = $8
		RETURNING 
			TO_CHAR(expense_date, 'YYYY-MM-DD'), created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		e.Category, e.Description, e.Amount, e.ExpenseDate,
		e.IsRecurring, e.RecurringPeriod, e.Notes, id,
	).Scan(&e.ExpenseDate, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, err
	}
	e.ID = id
	return e, nil
}

// Delete manghapus data expense
func (r *ExpenseRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM expenses WHERE id = $1`, id)
	return err
}

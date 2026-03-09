package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type LoyaltyRepository struct {
	db *sql.DB
}

func NewLoyaltyRepository(db *sql.DB) *LoyaltyRepository {
	return &LoyaltyRepository{db: db}
}

// GetHistoryByCustomerID retrieves the point history for a specific customer
func (r *LoyaltyRepository) GetHistoryByCustomerID(customerID int) ([]models.LoyaltyTransaction, error) {
	query := `
		SELECT id, customer_id, transaction_id, type, points, description, created_by_user_id, created_at
		FROM loyalty_transactions
		WHERE customer_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`
	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.LoyaltyTransaction
	for rows.Next() {
		var lt models.LoyaltyTransaction
		var txID sql.NullInt64
		var desc sql.NullString
		var createdBy sql.NullInt64
		err := rows.Scan(
			&lt.ID, &lt.CustomerID, &txID, &lt.Type, &lt.Points, &desc, &createdBy, &lt.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if txID.Valid {
			id := int(txID.Int64)
			lt.TransactionID = &id
		}
		if desc.Valid {
			d := desc.String
			lt.Description = &d
		}
		if createdBy.Valid {
			id := int(createdBy.Int64)
			lt.CreatedBy = &id
		}

		history = append(history, lt)
	}
	return history, nil
}

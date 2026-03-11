package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type SubscriptionPackageRepository struct {
	db *sql.DB
}

func NewSubscriptionPackageRepository(db *sql.DB) *SubscriptionPackageRepository {
	return &SubscriptionPackageRepository{db: db}
}

func (r *SubscriptionPackageRepository) GetAll() ([]models.SubscriptionPackage, error) {
	query := `SELECT id, name, max_kasir, max_products, price, features, is_active, created_at FROM subscription_packages ORDER BY id ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pkgs []models.SubscriptionPackage
	for rows.Next() {
		var p models.SubscriptionPackage
		if err := rows.Scan(
			&p.ID, &p.Name, &p.MaxKasir, &p.MaxProducts, &p.Price, &p.Features, &p.IsActive, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		pkgs = append(pkgs, p)
	}
	return pkgs, nil
}

// Tambahkan func Create, Update, GetByID nnt jika diperlukan untuk Superadmin

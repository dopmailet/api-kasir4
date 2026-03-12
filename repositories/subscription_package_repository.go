package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type SubscriptionPackageRepository struct {
	db *sql.DB
}

func NewSubscriptionPackageRepository(db *sql.DB) *SubscriptionPackageRepository {
	return &SubscriptionPackageRepository{db: db}
}

func (r *SubscriptionPackageRepository) GetAll(publicOnly bool) ([]models.SubscriptionPackage, error) {
	query := `SELECT id, name, max_kasir, max_products, price, is_active,
       description, features, period, discount_percent, discount_label, is_popular, sort_order,
       created_at, updated_at
FROM subscription_packages`

	if publicOnly {
		query += ` WHERE is_active = true ORDER BY sort_order ASC, price ASC`
	} else {
		query += ` ORDER BY sort_order ASC, id ASC`
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pkgs []models.SubscriptionPackage
	for rows.Next() {
		var p models.SubscriptionPackage
		if err := rows.Scan(
			&p.ID, &p.Name, &p.MaxKasir, &p.MaxProducts, &p.Price, &p.IsActive,
			&p.Description, &p.Features, &p.Period, &p.DiscountPercent, &p.DiscountLabel,
			&p.IsPopular, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pkgs = append(pkgs, p)
	}
	return pkgs, nil
}

// GetByID mengambil satu paket berdasarkan ID
func (r *SubscriptionPackageRepository) GetByID(id int) (*models.SubscriptionPackage, error) {
	query := `SELECT id, name, max_kasir, max_products, price, is_active,
       description, features, period, discount_percent, discount_label, is_popular, sort_order,
       created_at, updated_at
FROM subscription_packages WHERE id = $1`

	var p models.SubscriptionPackage
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.MaxKasir, &p.MaxProducts, &p.Price, &p.IsActive,
		&p.Description, &p.Features, &p.Period, &p.DiscountPercent, &p.DiscountLabel,
		&p.IsPopular, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Create menambahkan paket langganan baru
func (r *SubscriptionPackageRepository) Create(p *models.SubscriptionPackage) error {
	query := `INSERT INTO subscription_packages 
  (name, max_kasir, max_products, price, is_active, 
   description, features, period, discount_percent, discount_label, is_popular, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 
        COALESCE((SELECT MAX(sort_order) FROM subscription_packages), 0) + 1)
RETURNING id, sort_order, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		p.Name, p.MaxKasir, p.MaxProducts, p.Price, p.IsActive,
		p.Description, p.Features, p.Period, p.DiscountPercent, p.DiscountLabel, p.IsPopular,
	).Scan(&p.ID, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt)

	return err
}

// Update memodifikasi paket yang ada
func (r *SubscriptionPackageRepository) Update(p *models.SubscriptionPackage) error {
	// Jika diset priority, matikan priority paket lain dulu
	if p.IsPopular {
		_, _ = r.db.Exec(`UPDATE subscription_packages SET is_popular = false WHERE id != $1`, p.ID)
	}

	query := `UPDATE subscription_packages
SET name = COALESCE(NULLIF($2, ''), name),
    max_kasir = COALESCE(NULLIF($3, 0), max_kasir),
    max_products = COALESCE(NULLIF($4, 0), max_products),
    price = $5,
    is_active = $6,
    description = COALESCE(NULLIF($7, ''), description),
    features = COALESCE($8, features),
    period = COALESCE(NULLIF($9, ''), period),
    discount_percent = $10,
    discount_label = $11,
    is_popular = $12,
    updated_at = NOW()
WHERE id = $1
RETURNING sort_order, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		p.ID, p.Name, p.MaxKasir, p.MaxProducts, p.Price, p.IsActive,
		p.Description, p.Features, p.Period, p.DiscountPercent, p.DiscountLabel, p.IsPopular,
	).Scan(&p.SortOrder, &p.CreatedAt, &p.UpdatedAt)

	return err
}

// TogglePopular mengubah status populer satu paket sekaligus mereset paket lain
// Endpoint ini TIDAK membutuhkan seluruh data paket - hanya butuh ID dan nilai is_popular
func (r *SubscriptionPackageRepository) TogglePopular(id int, isPopular bool) error {
	// Jika mengaktifkan, matikan flag popular di semua paket lain dulu
	if isPopular {
		_, err := r.db.Exec(`UPDATE subscription_packages SET is_popular = false WHERE id != $1`, id)
		if err != nil {
			return err
		}
	}

	// Update hanya field is_popular di paket ini
	_, err := r.db.Exec(
		`UPDATE subscription_packages SET is_popular = $2, updated_at = NOW() WHERE id = $1`,
		id, isPopular,
	)
	return err
}

// Delete menonaktifkan atau menghapus layanan (jika tidak dipakai)
func (r *SubscriptionPackageRepository) Delete(id int) error {
	// Menilai Apakah Ada store yang masih memakainya
	var activeCount int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM stores WHERE subscription_package_id = $1 AND (subscription_end_date IS NULL OR subscription_end_date > NOW())`, id).Scan(&activeCount)
	if err != nil {
		return err
	}

	if activeCount > 0 {
		return fmt.Errorf("Tidak dapat menghapus paket yang masih digunakan oleh %d toko aktif", activeCount)
	}

	// Reset store yang expired
	_, _ = r.db.Exec(`UPDATE stores SET subscription_package_id = NULL WHERE subscription_package_id = $1`, id)

	// Hapus paket
	_, err = r.db.Exec(`DELETE FROM subscription_packages WHERE id = $1`, id)
	return err
}

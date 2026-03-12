package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"log"
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
       max_daily_sales, created_at, updated_at
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
			&p.IsPopular, &p.SortOrder, &p.MaxDailySales, &p.CreatedAt, &p.UpdatedAt,
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
       max_daily_sales, created_at, updated_at
FROM subscription_packages WHERE id = $1`

	var p models.SubscriptionPackage
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.MaxKasir, &p.MaxProducts, &p.Price, &p.IsActive,
		&p.Description, &p.Features, &p.Period, &p.DiscountPercent, &p.DiscountLabel,
		&p.IsPopular, &p.SortOrder, &p.MaxDailySales, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Create menambahkan paket langganan baru.
// featuresJSON harus berupa JSON string valid, mis: "[]" atau '["Fitur A","Fitur B"]'
func (r *SubscriptionPackageRepository) Create(p *models.SubscriptionPackage, featuresJSON string) error {
	query := `INSERT INTO subscription_packages
  (name, max_kasir, max_products, price, is_active,
   description, features, period, discount_percent, discount_label, is_popular, max_daily_sales, sort_order)
VALUES ($1, $2, $3, $4, $5, NULLIF(TRIM($6), ''), $7::jsonb, $8, $9, NULLIF(TRIM($10), ''), $11, $12,
        COALESCE((SELECT MAX(sort_order) FROM subscription_packages), 0) + 1)
RETURNING id, sort_order, created_at, updated_at`

	// Ekstrak nilai nullable dari pointer — kirim "" jika nil agar NULLIF bekerja
	descVal := ""
	if p.Description != nil {
		descVal = *p.Description
	}
	discountLabelVal := ""
	if p.DiscountLabel != nil {
		discountLabelVal = *p.DiscountLabel
	}

	err := r.db.QueryRow(
		query,
		p.Name, p.MaxKasir, p.MaxProducts, p.Price, p.IsActive,
		descVal, featuresJSON, p.Period, p.DiscountPercent, discountLabelVal, p.IsPopular, p.MaxDailySales,
	).Scan(&p.ID, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		log.Printf("[Create Package] SQL error: %v | featuresJSON=%s", err, featuresJSON)
		return err
	}
	return nil
}

// Update memodifikasi paket yang ada secara atomik.
// featuresJSON harus berupa JSON string valid, mis: "[]" atau '["Fitur A","Fitur B"]'
// Jika is_popular = true, paket lain di-reset false dalam satu transaksi.
func (r *SubscriptionPackageRepository) Update(p *models.SubscriptionPackage, featuresJSON string) error {
	// Ekstrak nilai nullable dari pointer — kirim "" jika nil agar NULLIF bekerja
	descVal := ""
	if p.Description != nil {
		descVal = *p.Description
	}
	discountLabelVal := ""
	if p.DiscountLabel != nil {
		discountLabelVal = *p.DiscountLabel
	}

	// Gunakan transaksi agar reset is_popular + update data atomik (tidak ada race condition)
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("[Update Package %d] Gagal begin transaction: %v", p.ID, err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Jika paket ini akan jadi populer — reset semua paket lain dulu
	if p.IsPopular {
		_, err = tx.Exec(
			`UPDATE subscription_packages SET is_popular = false, updated_at = NOW() WHERE id != $1`,
			p.ID,
		)
		if err != nil {
			log.Printf("[Update Package %d] Gagal reset is_popular paket lain: %v", p.ID, err)
			return err
		}
	}

	query := `UPDATE subscription_packages
SET name            = $2,
    max_kasir       = $3,
    max_products    = $4,
    price           = $5,
    is_active       = $6,
    description     = NULLIF(TRIM($7), ''),
    features        = $8::jsonb,
    period          = $9,
    discount_percent= $10,
    discount_label  = NULLIF(TRIM($11), ''),
    is_popular      = $12,
    max_daily_sales = $13,
    updated_at      = NOW()
WHERE id = $1
RETURNING sort_order, created_at, updated_at`

	err = tx.QueryRow(
		query,
		p.ID, p.Name, p.MaxKasir, p.MaxProducts, p.Price, p.IsActive,
		descVal, featuresJSON, p.Period, p.DiscountPercent, discountLabelVal, p.IsPopular, p.MaxDailySales,
	).Scan(&p.SortOrder, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[Update Package %d] Not found", p.ID)
			return fmt.Errorf("package not found")
		}
		log.Printf("[Update Package %d] SQL error: %v | featuresJSON=%s | desc=%q | discountLabel=%q",
			p.ID, err, featuresJSON, descVal, discountLabelVal)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[Update Package %d] Gagal commit transaction: %v", p.ID, err)
		return err
	}

	return nil
}

// TogglePopular mengubah status populer satu paket sekaligus mereset paket lain secara atomik.
// Endpoint ini TIDAK membutuhkan seluruh data paket - hanya butuh ID dan nilai is_popular.
func (r *SubscriptionPackageRepository) TogglePopular(id int, isPopular bool) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Jika mengaktifkan, matikan flag popular di semua paket lain dulu
	if isPopular {
		_, err = tx.Exec(
			`UPDATE subscription_packages SET is_popular = false, updated_at = NOW() WHERE id != $1`,
			id,
		)
		if err != nil {
			log.Printf("[TogglePopular %d] Gagal reset paket lain: %v", id, err)
			return err
		}
	}

	// Update hanya field is_popular di paket ini
	_, err = tx.Exec(
		`UPDATE subscription_packages SET is_popular = $2, updated_at = NOW() WHERE id = $1`,
		id, isPopular,
	)
	if err != nil {
		log.Printf("[TogglePopular %d] Gagal update: %v", id, err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[TogglePopular %d] Gagal commit: %v", id, err)
		return err
	}

	return nil
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

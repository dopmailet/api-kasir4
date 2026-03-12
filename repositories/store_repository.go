package repositories

import (
	"database/sql"
	"kasir-api/models"
	"time"
)

type StoreRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

// Create new store and assign default package (Free Trial/Gratis)
func (r *StoreRepository) Create(store *models.Store, tx *sql.Tx) error {
	query := `
		INSERT INTO stores (name, address, phone, email, subscription_package_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	// Assuming subscription_package_id = 1 is the default Free Package
	// In production, we should dynamically fetch it or make it configurable
	store.SubscriptionPackageID = 1
	store.IsActive = true

	var err error
	if tx != nil {
		err = tx.QueryRow(query, store.Name, store.Address, store.Phone, store.Email, store.SubscriptionPackageID, store.IsActive).Scan(&store.ID, &store.CreatedAt, &store.UpdatedAt)
	} else {
		err = r.db.QueryRow(query, store.Name, store.Address, store.Phone, store.Email, store.SubscriptionPackageID, store.IsActive).Scan(&store.ID, &store.CreatedAt, &store.UpdatedAt)
	}

	return err
}

func (r *StoreRepository) GetByID(id int) (*models.Store, error) {
	query := `SELECT id, name, address, phone, email, subscription_package_id, subscription_end_date, is_active, is_verified, created_at, updated_at FROM stores WHERE id = $1`
	var store models.Store
	err := r.db.QueryRow(query, id).Scan(
		&store.ID, &store.Name, &store.Address, &store.Phone, &store.Email,
		&store.SubscriptionPackageID, &store.SubscriptionEndDate, &store.IsActive, &store.IsVerified,
		&store.CreatedAt, &store.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *StoreRepository) GetAll() ([]models.Store, error) {
	query := `SELECT id, name, address, phone, email, subscription_package_id, subscription_end_date, is_active, is_verified, created_at, updated_at FROM stores ORDER BY id ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stores []models.Store
	for rows.Next() {
		var s models.Store
		if err := rows.Scan(
			&s.ID, &s.Name, &s.Address, &s.Phone, &s.Email,
			&s.SubscriptionPackageID, &s.SubscriptionEndDate, &s.IsActive, &s.IsVerified,
			&s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		stores = append(stores, s)
	}
	return stores, nil
}

func (r *StoreRepository) Update(store *models.Store) error {
	query := `
		UPDATE stores 
		SET name = $1, address = $2, phone = $3, email = $4, 
		    subscription_package_id = $5, subscription_end_date = $6, is_active = $7
		WHERE id = $8
	`
	_, err := r.db.Exec(query,
		store.Name, store.Address, store.Phone, store.Email,
		store.SubscriptionPackageID, store.SubscriptionEndDate, store.IsActive, store.ID,
	)
	return err
}

// UpdatePackage mengganti paket langganan sebuah toko beserta tanggal berlakunya
// Jika packageID > 1 (paket berbayar), otomatis set is_verified = true (permanen)
func (r *StoreRepository) UpdatePackage(storeID int, packageID int, endDate *time.Time) error {
	query := `
		UPDATE stores
		SET subscription_package_id = $1, subscription_end_date = $2,
		    is_verified = CASE WHEN $1 > 1 THEN true ELSE is_verified END
		WHERE id = $3
	`
	_, err := r.db.Exec(query, packageID, endDate, storeID)
	return err
}

// Delete menghapus satu toko berdasarkan ID (hard delete)
func (r *StoreRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM stores WHERE id = $1`, id)
	return err
}

// DeleteUnverified menghapus semua toko yang belum pernah berlangganan berbayar
func (r *StoreRepository) DeleteUnverified() (int64, error) {
	res, err := r.db.Exec(`DELETE FROM stores WHERE is_verified = false AND id != 1`)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

// DeleteExpired menghapus toko yang tidak terverifikasi DAN langganannya sudah habis
// (is_verified=false artinya belum pernah berlangganan berbayar sama sekali)
func (r *StoreRepository) DeleteNeverSubscribed() (int64, error) {
	res, err := r.db.Exec(`
		DELETE FROM stores 
		WHERE is_verified = false 
		  AND subscription_package_id = 1 
		  AND id != 1
	`)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

// CountActiveCashiers menghitung jumlah akun kasir yang aktif dalam satu toko
func (r *StoreRepository) CountActiveCashiers(storeID int) (int, error) {
	var count int
	// Hanya menghitung role kasir yang aktif
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE store_id = $1 AND role = 'kasir' AND is_active = true`, storeID).Scan(&count)
	return count, err
}

// CountActiveProducts menghitung jumlah produk untuk sebuah toko
func (r *StoreRepository) CountActiveProducts(storeID int) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM products WHERE store_id = $1`, storeID).Scan(&count)
	return count, err
}

// CountTodayTransactions menghitung jumlah transaksi (checkout) hari ini untuk sebuah toko
func (r *StoreRepository) CountTodayTransactions(storeID int, timezone string) (int, error) {
	if timezone == "" {
		timezone = "Asia/Makassar"
	}

	var count int
	query := `
		SELECT COUNT(*) 
		FROM transactions 
		WHERE store_id = $1 
		  AND created_at >= (CURRENT_DATE AT TIME ZONE $2)
		  AND created_at < (CURRENT_DATE AT TIME ZONE $2) + INTERVAL '1 day'
	`
	err := r.db.QueryRow(query, storeID, timezone).Scan(&count)
	return count, err
}

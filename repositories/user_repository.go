package repositories

import (
	"database/sql"
	"kasir-api/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByUsername retrieves a user by username
// Digunakan untuk login
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, password, nama_lengkap, role, is_active,
		       COALESCE(store_id, 0), COALESCE(is_superadmin, false),
		       created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.NamaLengkap,
		&user.Role,
		&user.IsActive,
		&user.StoreID,
		&user.IsSuperadmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, username, password, nama_lengkap, role, is_active,
		       COALESCE(store_id, 0), COALESCE(is_superadmin, false),
		       created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.NamaLengkap,
		&user.Role,
		&user.IsActive,
		&user.StoreID,
		&user.IsSuperadmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAll retrieves all users (untuk admin)
func (r *UserRepository) GetAll(storeID int) ([]models.User, error) {
	query := `
		SELECT id, username, nama_lengkap, role, is_active, created_at, updated_at
		FROM users
		WHERE store_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.Query(query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.NamaLengkap,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Create adds a new user
func (r *UserRepository) Create(user *models.User, tx ...*sql.Tx) error {
	query := `
		INSERT INTO users (username, password, nama_lengkap, role, is_active, store_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	var err error
	if len(tx) > 0 && tx[0] != nil {
		err = tx[0].QueryRow(
			query,
			user.Username,
			user.Password,
			user.NamaLengkap,
			user.Role,
			user.IsActive,
			user.StoreID,
		).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	} else {
		err = r.db.QueryRow(
			query,
			user.Username,
			user.Password,
			user.NamaLengkap,
			user.Role,
			user.IsActive,
			user.StoreID,
		).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	}

	return err
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET nama_lengkap = $1, role = $2, is_active = $3
		WHERE id = $4 AND store_id = $5
	`

	_, err := r.db.Exec(query, user.NamaLengkap, user.Role, user.IsActive, user.ID, user.StoreID)
	return err
}

// UpdatePassword updates user password
func (r *UserRepository) UpdatePassword(userID int, hashedPassword string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := r.db.Exec(query, hashedPassword, userID)
	return err
}

// Delete removes a user (soft delete dengan set is_active = false lebih baik)
func (r *UserRepository) Delete(id int, storeID int) error {
	query := `UPDATE users SET is_active = false WHERE id = $1 AND store_id = $2`
	_, err := r.db.Exec(query, id, storeID)
	return err
}

// SetActive mengaktifkan/menonaktifkan user
func (r *UserRepository) SetActive(id int, isActive bool, storeID int) error {
	query := `UPDATE users SET is_active = $1 WHERE id = $2 AND store_id = $3`
	_, err := r.db.Exec(query, isActive, id, storeID)
	return err
}

// GetAllStoreUsers retrieves all users (admin and kasir) from all stores for superadmin.
func (r *UserRepository) GetAllStoreUsers() ([]models.StoreUserResponse, error) {
	query := `
		SELECT 
			u.id, 
			u.store_id, 
			u.username, 
			s.email, 
			u.role, 
			u.nama_lengkap as full_name, 
			u.created_at
		FROM users u
		LEFT JOIN stores s ON u.store_id = s.id
		WHERE u.role IN ('admin', 'kasir')
		ORDER BY u.store_id ASC, u.role ASC, u.created_at ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.StoreUserResponse
	for rows.Next() {
		var u models.StoreUserResponse
		err := rows.Scan(&u.ID, &u.StoreID, &u.Username, &u.Email, &u.Role, &u.FullName, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

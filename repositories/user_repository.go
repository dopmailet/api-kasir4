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
		SELECT id, username, password, nama_lengkap, role, is_active, created_at, updated_at
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
		SELECT id, username, password, nama_lengkap, role, is_active, created_at, updated_at
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
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAll retrieves all users (untuk admin)
func (r *UserRepository) GetAll() ([]models.User, error) {
	query := `
		SELECT id, username, nama_lengkap, role, is_active, created_at, updated_at
		FROM users
		ORDER BY id ASC
	`

	rows, err := r.db.Query(query)
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
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, password, nama_lengkap, role, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Password,
		user.NamaLengkap,
		user.Role,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET nama_lengkap = $1, role = $2, is_active = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(query, user.NamaLengkap, user.Role, user.IsActive, user.ID)
	return err
}

// UpdatePassword updates user password
func (r *UserRepository) UpdatePassword(userID int, hashedPassword string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := r.db.Exec(query, hashedPassword, userID)
	return err
}

// Delete removes a user (soft delete dengan set is_active = false lebih baik)
func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// SetActive mengaktifkan/menonaktifkan user
func (r *UserRepository) SetActive(id int, isActive bool) error {
	query := `UPDATE users SET is_active = $1 WHERE id = $2`
	_, err := r.db.Exec(query, isActive, id)
	return err
}

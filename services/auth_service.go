package services

import (
	"database/sql"
	"kasir-api/models"
	"kasir-api/repositories"
	"kasir-api/utils"
)

// AuthService handles authentication logic
type AuthService struct {
	db        *sql.DB
	userRepo  *repositories.UserRepository
	storeRepo *repositories.StoreRepository
}

// NewAuthService creates a new AuthService
func NewAuthService(db *sql.DB, userRepo *repositories.UserRepository, storeRepo *repositories.StoreRepository) *AuthService {
	return &AuthService{
		db:        db,
		userRepo:  userRepo,
		storeRepo: storeRepo,
	}
}

// Login melakukan autentikasi user dan mengembalikan JWT token
func (s *AuthService) Login(username, password string) (*models.LoginResponse, error) {
	// 1. Cari user berdasarkan username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrInvalidCredentials
		}
		return nil, err
	}

	// 2. Cek apakah user aktif
	if !user.IsActive {
		return nil, models.ErrUserInactive
	}

	// 3. Validasi password
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, models.ErrInvalidCredentials
	}

	// 4. Generate JWT token (sertakan store_id dan is_superadmin)
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role, user.StoreID, user.IsSuperadmin)
	if err != nil {
		return nil, err
	}

	// 5. Return response (password tidak di-include karena json:"-")
	return &models.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Register membuat user baru (hanya bisa dilakukan oleh admin)
func (s *AuthService) Register(username, password, namaLengkap, role string) (*models.User, error) {
	// 1. Validasi role
	if role != models.RoleAdmin && role != models.RoleKasir {
		return nil, models.ErrInvalidRole
	}

	// 2. Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 3. Buat user baru
	user := &models.User{
		Username:    username,
		Password:    hashedPassword,
		NamaLengkap: namaLengkap,
		Role:        role,
		IsActive:    true,
	}

	// 4. Simpan ke database
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// RegisterStore makes a new tenant along with the first admin user
func (s *AuthService) RegisterStore(req models.StoreRegisterRequest) (*models.StoreRegisterResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Rollback jika ada yang panic atau err di tengah jalan

	// 1. Create Store
	store := &models.Store{
		Name:    req.StoreName,
		Email:   &req.AdminEmail, // Gunakan admin email sbg email toko untuk default
		Address: nil,
		Phone:   nil,
	}
	if err := s.storeRepo.Create(store, tx); err != nil {
		return nil, err
	}

	// 2. Hash Password for Admin
	hashedPassword, err := utils.HashPassword(req.AdminPassword)
	if err != nil {
		return nil, err
	}

	// 3. Create Admin User yang terikat pada Toko yang baru dibuat
	userAdmin := &models.User{
		Username:    req.AdminUsername,
		Password:    hashedPassword,
		NamaLengkap: req.AdminName,
		Role:        models.RoleAdmin, // Auto-assign Admin
		StoreID:     store.ID,
		IsActive:    true,
	}

	if err := s.userRepo.Create(userAdmin, tx); err != nil {
		return nil, err
	}

	// 4. Commit Transaksi
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// 5. Generate JWT untuk Auto-Login
	token, err := utils.GenerateJWT(userAdmin.ID, userAdmin.Username, userAdmin.Role, userAdmin.StoreID, userAdmin.IsSuperadmin)
	if err != nil {
		return nil, err
	}

	return &models.StoreRegisterResponse{
		Message: "Pendaftaran toko berhasil",
		Store:   store,
		Token:   token,
	}, nil
}

// ChangePassword mengubah password user
func (s *AuthService) ChangePassword(userID int, oldPassword, newPassword string) error {
	// 1. Ambil user dari database
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return models.ErrUserNotFound
	}

	// 2. Validasi old password
	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return models.ErrInvalidCredentials
	}

	// 3. Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 4. Update password di database
	return s.userRepo.UpdatePassword(userID, hashedPassword)
}

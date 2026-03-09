package services

import (
	"database/sql"
	"kasir-api/models"
	"kasir-api/repositories"
	"kasir-api/utils"
)

// AuthService handles authentication logic
type AuthService struct {
	userRepo *repositories.UserRepository
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
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

	// 4. Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
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

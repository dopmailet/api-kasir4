package services

import (
	"fmt"
	"kasir-api/models"
	"kasir-api/repositories"
	"kasir-api/utils"

	"golang.org/x/crypto/bcrypt"
)

// UserService handles user management logic
type UserService struct {
	userRepo *repositories.UserRepository
	storeSvc *StoreService
}

// NewUserService creates a new UserService
func NewUserService(userRepo *repositories.UserRepository, storeSvc *StoreService) *UserService {
	return &UserService{
		userRepo: userRepo,
		storeSvc: storeSvc,
	}
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers(storeID int) ([]models.User, error) {
	return s.userRepo.GetAll(storeID)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username, password, role string, storeID int) (*models.User, error) {
	if role != models.RoleKasir {
		return nil, fmt.Errorf("hanya role 'kasir' yang dapat dibuat")
	}

	// Cek batasan paket berlangganan untuk kasir
	limits, err := s.storeSvc.GetStoreLimits(storeID, "Asia/Makassar")
	if err != nil {
		return nil, err
	}

	if limits.CurrentCashiers >= limits.MaxCashiers {
		return nil, fmt.Errorf("Batas kasir untuk paket %s adalah %d. Silakan upgrade paket.", limits.PackageName, limits.MaxCashiers)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:    username,
		Password:    hashedPassword,
		NamaLengkap: username, // Default to username if not provided
		Role:        role,
		IsActive:    true,
		StoreID:     storeID,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdatePassword updates a user's password (by admin)
// Requires the admin's own current password to authorize the change.
func (s *UserService) UpdatePassword(targetID int, newPassword string, currentAdminID int, currentPassword string, storeID int) error {
	// 1. Verifikasi user target ada
	targetUser, err := s.userRepo.GetByID(targetID)
	if err != nil {
		return err
	}
	if targetUser.StoreID != storeID {
		return models.ErrUserNotFound // mencegah edit lintas toko
	}

	// 2. Ambil password hash admin yang sedang login
	admin, err := s.userRepo.GetByID(currentAdminID)
	if err != nil {
		return models.ErrUserNotFound
	}

	// 3. Bandingkan current_password dengan hash admin
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(currentPassword)); err != nil {
		return models.ErrInvalidCredentials
	}

	// 4. Hash password baru
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(targetID, hashedPassword)
}

// DeleteUser deletes a user (except themselves)
func (s *UserService) DeleteUser(targetID, currentAdminID int, storeID int) error {
	if targetID == currentAdminID {
		return models.ErrCannotDeleteSelf
	}

	targetUser, err := s.userRepo.GetByID(targetID)
	if err != nil {
		return err // Check if user exists before deleting
	}
	if targetUser.StoreID != storeID {
		return models.ErrUserNotFound // mencegah lintas toko
	}

	return s.userRepo.Delete(targetID, storeID)
}

// UpdateStatus updates the active status of a user
func (s *UserService) UpdateStatus(targetID int, isActive bool, currentAdminID int, storeID int) (*models.User, error) {
	if targetID == currentAdminID {
		return nil, fmt.Errorf("anda tidak dapat menonaktifkan akun sendiri")
	}

	targetUser, err := s.userRepo.GetByID(targetID)
	if err != nil {
		return nil, models.ErrUserNotFound
	}
	if targetUser.StoreID != storeID {
		return nil, models.ErrUserNotFound
	}

	if targetUser.Role == models.RoleAdmin && !isActive {
		return nil, fmt.Errorf("user dengan role admin tidak dapat dinonaktifkan")
	}

	err = s.userRepo.SetActive(targetID, isActive, storeID)
	if err != nil {
		return nil, err
	}

	updatedUser, _ := s.userRepo.GetByID(targetID)
	return updatedUser, nil
}

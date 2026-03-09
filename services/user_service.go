package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"kasir-api/utils"

	"golang.org/x/crypto/bcrypt"
)

// UserService handles user management logic
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username, password, role string) (*models.User, error) {
	if role != models.RoleAdmin && role != models.RoleKasir {
		return nil, models.ErrInvalidRole
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
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdatePassword updates a user's password (by admin)
// Requires the admin's own current password to authorize the change.
func (s *UserService) UpdatePassword(targetID int, newPassword string, currentAdminID int, currentPassword string) error {
	// 1. Verifikasi user target ada
	_, err := s.userRepo.GetByID(targetID)
	if err != nil {
		return err
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
func (s *UserService) DeleteUser(targetID, currentAdminID int) error {
	if targetID == currentAdminID {
		return models.ErrCannotDeleteSelf
	}

	_, err := s.userRepo.GetByID(targetID)
	if err != nil {
		return err // Check if user exists before deleting
	}

	return s.userRepo.Delete(targetID)
}

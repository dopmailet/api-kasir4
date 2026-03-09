package services

import (
	"errors"
	"fmt"
	"kasir-api/models"       // Import models untuk struct Category
	"kasir-api/repositories" // Import repositories untuk akses database
	"log"
	"strings"
)

// CategoryService handles business logic for categories
// Service adalah layer antara Handler dan Repository
type CategoryService struct {
	repo  *repositories.CategoryRepository // Pointer ke CategoryRepository
	cache *CacheService                    // Pointer ke CacheService untuk Redis
}

// NewCategoryService creates a new CategoryService
// Fungsi ini adalah "constructor" untuk membuat instance CategoryService
// Sekarang menerima CacheService untuk caching
func NewCategoryService(repo *repositories.CategoryRepository, cache *CacheService) *CategoryService {
	return &CategoryService{
		repo:  repo,
		cache: cache,
	}
}

// validateCategory melakukan validasi data kategori
// Fungsi ini akan dipanggil sebelum Create dan Update
func (s *CategoryService) validateCategory(category *models.Category) error {
	// Validasi nama tidak boleh kosong
	if strings.TrimSpace(category.Nama) == "" {
		return errors.New("nama kategori tidak boleh kosong")
	}

	// Validasi nama minimal 2 karakter
	if len(strings.TrimSpace(category.Nama)) < 2 {
		return errors.New("nama kategori minimal 2 karakter")
	}

	return nil
}

// GetAll retrieves all categories with caching
// Fungsi ini memanggil repository untuk ambil semua kategori
func (s *CategoryService) GetAll() ([]models.Category, error) {
	// Generate cache key
	cacheKey := s.cache.GenerateKey("categories", "list", "all")

	// Coba ambil dari cache
	var categories []models.Category
	if s.cache.Get(cacheKey, &categories) {
		// Cache HIT - return dari cache
		return categories, nil
	}

	// Cache MISS - ambil dari database
	categories, err := s.repo.GetAll()
	if err != nil {
		log.Printf("❌ Error getting categories from database: %v", err)
		return nil, err
	}

	// Simpan ke cache
	s.cache.Set(cacheKey, categories, 0) // 0 = gunakan default TTL (5 menit)

	return categories, nil
}

// GetByID retrieves a category by ID with caching
// Fungsi ini memanggil repository untuk ambil 1 kategori by ID
func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	// Generate cache key
	cacheKey := s.cache.GenerateKey("categories", "detail", fmt.Sprintf("id:%d", id))

	// Coba ambil dari cache
	var category models.Category
	if s.cache.Get(cacheKey, &category) {
		// Cache HIT - return dari cache
		return &category, nil
	}

	// Cache MISS - ambil dari database
	categoryPtr, err := s.repo.GetByID(id)
	if err != nil {
		log.Printf("❌ Error getting category by ID %d: %v", id, err)
		return nil, err
	}

	// Simpan ke cache
	s.cache.Set(cacheKey, categoryPtr, 0)

	return categoryPtr, nil
}

// Create adds a new category and invalidates cache
// Fungsi ini memanggil repository untuk tambah kategori baru
func (s *CategoryService) Create(category *models.Category) error {
	// Validasi input
	if err := s.validateCategory(category); err != nil {
		log.Printf("⚠️ Validation error on create category: %v", err)
		return err
	}

	// Trim whitespace dari nama
	category.Nama = strings.TrimSpace(category.Nama)

	// Panggil repository untuk save ke database
	err := s.repo.Create(category)
	if err != nil {
		log.Printf("❌ Error creating category: %v", err)
		return err
	}

	log.Printf("✅ Category created successfully: ID=%d, Name=%s", category.ID, category.Nama)

	// Invalidate semua cache categories karena ada data baru
	s.cache.DeletePattern("categories:*")

	return nil
}

// Update updates an existing category and invalidates cache
// Fungsi ini memanggil repository untuk update kategori
func (s *CategoryService) Update(id int, category *models.Category) error {
	// Set ID untuk memastikan update kategori yang benar
	category.ID = id

	// Validasi input
	if err := s.validateCategory(category); err != nil {
		log.Printf("⚠️ Validation error on update category ID %d: %v", id, err)
		return err
	}

	// Trim whitespace dari nama
	category.Nama = strings.TrimSpace(category.Nama)

	// Panggil repository untuk update di database
	err := s.repo.Update(category)
	if err != nil {
		log.Printf("❌ Error updating category ID %d: %v", id, err)
		return err
	}

	log.Printf("✅ Category updated successfully: ID=%d, Name=%s", id, category.Nama)

	// Invalidate cache untuk kategori ini dan semua list
	s.cache.Delete(s.cache.GenerateKey("categories", "detail", fmt.Sprintf("id:%d", id)))
	s.cache.DeletePattern("categories:list:*")

	return nil
}

// Delete removes a category and invalidates cache
// Fungsi ini memanggil repository untuk hapus kategori
func (s *CategoryService) Delete(id int) error {
	// Panggil repository untuk delete dari database
	err := s.repo.Delete(id)
	if err != nil {
		log.Printf("❌ Error deleting category ID %d: %v", id, err)
		return err
	}

	log.Printf("✅ Category deleted successfully: ID=%d", id)

	// Invalidate cache untuk kategori ini dan semua list
	s.cache.Delete(s.cache.GenerateKey("categories", "detail", fmt.Sprintf("id:%d", id)))
	s.cache.DeletePattern("categories:list:*")

	return nil
}

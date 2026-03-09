package services

import (
	"errors"
	"fmt"
	"kasir-api/models"       // Import models untuk struct Product
	"kasir-api/repositories" // Import repositories untuk akses database
	"log"
	"strings"
)

// ProductService handles business logic for products
// Service adalah layer antara Handler dan Repository
// Di sini kita bisa tambahkan validasi, business rules, dll
type ProductService struct {
	repo  *repositories.ProductRepository // Pointer ke ProductRepository
	cache *CacheService                   // Pointer ke CacheService untuk Redis
}

// NewProductService creates a new ProductService
// Fungsi ini adalah "constructor" untuk membuat instance ProductService
func NewProductService(repo *repositories.ProductRepository, cache *CacheService) *ProductService {
	return &ProductService{
		repo:  repo,
		cache: cache,
	}
}

// validateProduct melakukan validasi data produk
// Fungsi ini akan dipanggil sebelum Create dan Update
func (s *ProductService) validateProduct(product *models.Product) error {
	// Validasi nama tidak boleh kosong
	if strings.TrimSpace(product.Nama) == "" {
		return errors.New("nama produk tidak boleh kosong")
	}

	// Validasi nama minimal 2 karakter
	if len(strings.TrimSpace(product.Nama)) < 2 {
		return errors.New("nama produk minimal 2 karakter")
	}

	// Validasi harga harus lebih dari 0
	if product.Harga <= 0 {
		return errors.New("harga harus lebih dari 0")
	}

	// Validasi stok tidak boleh negatif
	if product.Stok < 0 {
		return errors.New("stok tidak boleh negatif")
	}

	return nil
}

// GetAll retrieves all products with caching and pagination
// Fungsi ini memanggil repository untuk ambil produk dengan Redis caching
// Parameter searchName untuk filter by name (kosong = ambil semua)
// Parameter pagination untuk limit dan offset (nil = tanpa pagination)
// Return: products, total count, error
func (s *ProductService) GetAll(searchName string, searchBarcode string, pagination *models.PaginationParams) ([]models.Product, int, error) {
	// Generate cache key berdasarkan search dan pagination
	cacheKey := s.cache.GenerateKey("products", "list",
		fmt.Sprintf("search:%s", searchName),
		fmt.Sprintf("barcode:%s", searchBarcode),
		fmt.Sprintf("page:%d", pagination.Page),
		fmt.Sprintf("limit:%d", pagination.Limit))

	// Struct untuk cache (products + total count)
	type CachedData struct {
		Products   []models.Product
		TotalCount int
	}

	// Coba ambil dari cache
	var cached CachedData
	if s.cache.Get(cacheKey, &cached) {
		// Cache HIT - return dari cache
		return cached.Products, cached.TotalCount, nil
	}

	// Cache MISS - ambil dari database
	products, totalCount, err := s.repo.GetAll(searchName, searchBarcode, pagination)
	if err != nil {
		log.Printf("❌ Error getting products from database: %v", err)
		return nil, 0, err
	}

	// Simpan ke cache untuk request berikutnya
	cached = CachedData{
		Products:   products,
		TotalCount: totalCount,
	}
	s.cache.Set(cacheKey, cached, 0) // 0 = gunakan default TTL (5 menit)

	return products, totalCount, nil
}

// GetByID retrieves a product by ID with caching
// Fungsi ini memanggil repository untuk ambil 1 produk by ID
func (s *ProductService) GetByID(id int) (*models.Product, error) {
	// Generate cache key
	cacheKey := s.cache.GenerateKey("products", "detail", fmt.Sprintf("id:%d", id))

	// Coba ambil dari cache
	var product models.Product
	if s.cache.Get(cacheKey, &product) {
		// Cache HIT - return dari cache
		return &product, nil
	}

	// Cache MISS - ambil dari database
	productPtr, err := s.repo.GetByID(id)
	if err != nil {
		log.Printf("❌ Error getting product by ID %d: %v", id, err)
		return nil, err
	}

	// Simpan ke cache
	s.cache.Set(cacheKey, productPtr, 0)

	return productPtr, nil
}

// GetByBarcode retrieves a product by barcode
// Fungsi ini memanggil repository untuk ambil 1 produk by barcode
func (s *ProductService) GetByBarcode(barcode string) (*models.Product, error) {
	// Generate cache key
	cacheKey := s.cache.GenerateKey("products", "barcode", fmt.Sprintf("code:%s", barcode))

	// Coba ambil dari cache
	var product models.Product
	if s.cache.Get(cacheKey, &product) {
		return &product, nil
	}

	// Cache MISS - ambil dari database
	productPtr, err := s.repo.GetByBarcode(barcode)
	if err != nil {
		log.Printf("❌ Error getting product by barcode %s: %v", barcode, err)
		return nil, err
	}

	// Simpan ke cache
	s.cache.Set(cacheKey, productPtr, 0)

	return productPtr, nil
}

// Create adds a new product and invalidates cache
// Fungsi ini memanggil repository untuk tambah produk baru
func (s *ProductService) Create(product *models.Product) error {
	// Validasi input
	if err := s.validateProduct(product); err != nil {
		log.Printf("⚠️ Validation error on create product: %v", err)
		return err
	}

	// Trim whitespace dari nama
	product.Nama = strings.TrimSpace(product.Nama)

	// Panggil repository untuk save ke database
	err := s.repo.Create(product)
	if err != nil {
		log.Printf("❌ Error creating product: %v", err)
		return err
	}

	log.Printf("✅ Product created successfully: ID=%d, Name=%s", product.ID, product.Nama)

	// Invalidate semua cache products list karena ada data baru
	s.cache.DeletePattern("products:list:*")

	return nil
}

// Update updates an existing product and invalidates cache
// Fungsi ini memanggil repository untuk update produk
func (s *ProductService) Update(id int, product *models.Product) error {
	// Set ID untuk memastikan update produk yang benar
	product.ID = id

	// Validasi input
	if err := s.validateProduct(product); err != nil {
		log.Printf("⚠️ Validation error on update product ID %d: %v", id, err)
		return err
	}

	// Trim whitespace dari nama
	product.Nama = strings.TrimSpace(product.Nama)

	// Panggil repository untuk update di database
	err := s.repo.Update(product)
	if err != nil {
		log.Printf("❌ Error updating product ID %d: %v", id, err)
		return err
	}

	log.Printf("✅ Product updated successfully: ID=%d, Name=%s", id, product.Nama)

	// Invalidate cache untuk produk ini dan semua list
	s.cache.Delete(s.cache.GenerateKey("products", "detail", fmt.Sprintf("id:%d", id)))
	s.cache.DeletePattern("products:list:*")

	return nil
}

// Delete removes a product and invalidates cache
// Fungsi ini memanggil repository untuk hapus produk
func (s *ProductService) Delete(id int) error {
	// Panggil repository untuk delete dari database
	err := s.repo.Delete(id)
	if err != nil {
		log.Printf("❌ Error deleting product ID %d: %v", id, err)
		return err
	}

	log.Printf("✅ Product deleted successfully: ID=%d", id)

	// Invalidate cache untuk produk ini dan semua list
	s.cache.Delete(s.cache.GenerateKey("products", "detail", fmt.Sprintf("id:%d", id)))
	s.cache.DeletePattern("products:list:*")

	return nil
}

package handlers

import (
	"encoding/json" // Package untuk encode/decode JSON
	"kasir-api/middleware"
	"kasir-api/models"   // Import models untuk struct Product
	"kasir-api/services" // Import services untuk business logic
	"log"
	"net/http" // Package untuk HTTP server
	"strconv"  // Package untuk convert string ke int
	"strings"  // Package untuk manipulasi string
)

// ProductHandler handles HTTP requests for products
// Handler adalah layer yang berhadapan langsung dengan HTTP request/response
type ProductHandler struct {
	service *services.ProductService // Pointer ke ProductService
}

// NewProductHandler creates a new ProductHandler
// Fungsi ini adalah "constructor" untuk membuat instance ProductHandler
func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service} // Return struct dengan service yang sudah di-inject
}

// HandleProducts handles /api/produk (GET all only)
// Produk baru sekarang dibuat lewat modul Pembelian (POST /api/purchases)
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.GetAll(w, r)
	case "POST":
		// Produk baru sekarang lewat modul pembelian
		http.Error(w, "Untuk menambah produk baru, gunakan POST /api/purchases", http.StatusMethodNotAllowed)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleProductByID handles /api/produk/{id} or /api/produk/barcode/{code}
// Fungsi ini handle: GET (by ID), PUT (update), DELETE (hapus), atau GET by barcode
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	// Cek apakah ini route barcode: /api/produk/barcode/{code}
	if strings.HasPrefix(r.URL.Path, "/api/produk/barcode/") {
		if r.Method == "GET" {
			h.GetByBarcode(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Switch berdasarkan HTTP method
	switch r.Method {
	case "GET":
		h.GetByID(w, r) // Kalau GET, panggil GetByID
	case "PUT":
		h.Update(w, r) // Kalau PUT, panggil Update
	case "DELETE":
		h.Delete(w, r) // Kalau DELETE, panggil Delete
	default:
		// Kalau method lain, return error
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetAll retrieves all products with pagination
// Fungsi ini handle GET /api/produk
// Support query parameter: ?name=xxx untuk search by name
// Support pagination: ?page=1&limit=10
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Ambil query parameter 'name' dari URL
	// Contoh: /api/produk?name=te -> searchName = "te"
	searchName := r.URL.Query().Get("name")

	// Parse pagination parameters
	page := 1
	limit := 10

	// Parse page parameter
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit parameter
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Buat pagination params
	pagination := models.NewPaginationParams(page, limit)

	// Ambil query parameter 'barcode' dari URL (exact match)
	searchBarcode := r.URL.Query().Get("barcode")

	// Panggil service untuk ambil produk (dengan filter dan pagination)
	products, totalCount, err := h.service.GetAll(searchName, searchBarcode, &pagination)
	if err != nil {
		// Log error untuk debugging
		log.Printf("❌ Handler: Error getting products: %v", err)
		// Kalau error, kirim HTTP error 500 (Internal Server Error)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter sensitive data based on role
	// Ambil user dari context
	user := middleware.GetUserFromContext(r.Context())

	// Jika user BUKAN admin, sembunyikan harga_beli dan margin
	// Note: Jika user nil (public access), juga sembunyikan
	if user == nil || !user.IsAdmin() {
		for i := range products {
			products[i].HargaBeli = nil
			products[i].Margin = nil
			products[i].CreatedBy = nil
		}
	} else {
		// Jika Admin, hitung margin untuk setiap produk
		for i := range products {
			products[i].Margin = products[i].CalculateMargin()
		}
	}

	// Buat response dengan pagination metadata
	response := models.PaginatedResponse{
		Data: products,
		Pagination: models.PaginationMeta{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			TotalItems: totalCount,
			TotalPages: models.CalculateTotalPages(totalCount, pagination.Limit),
		},
	}

	// Set header Content-Type jadi application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode response jadi JSON dan kirim ke client
	json.NewEncoder(w).Encode(response)
}

// GetByID retrieves a product by ID
// Fungsi ini handle GET /api/produk/{id}
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID dari URL
	// Misal URL: /api/produk/5 -> idStr = "5"
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// Convert string "5" jadi integer 5
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Log error untuk debugging
		log.Printf("⚠️ Handler: Invalid product ID: %s", idStr)
		// Kalau gagal convert (misal: /api/produk/abc), return error 400 (Bad Request)
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	// Panggil service untuk ambil produk by ID
	product, err := h.service.GetByID(id)
	if err != nil {
		// Log error untuk debugging
		log.Printf("❌ Handler: Error getting product ID %d: %v", id, err)
		// Kalau tidak ketemu, return error 404 (Not Found)
		http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
		return
	}

	// Filter sensitive data based on role
	user := middleware.GetUserFromContext(r.Context())

	// Jika user BUKAN admin, sembunyikan harga_beli dan margin
	if user == nil || !user.IsAdmin() {
		product.HargaBeli = nil
		product.Margin = nil
		product.CreatedBy = nil
	} else {
		// Jika Admin, hitung margin
		product.Margin = product.CalculateMargin()
	}

	// Set header Content-Type jadi application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode product jadi JSON dan kirim ke client
	json.NewEncoder(w).Encode(product)
}

// GetByBarcode retrieves a product by barcode
// Fungsi ini handle GET /api/produk/barcode/{code}
func (h *ProductHandler) GetByBarcode(w http.ResponseWriter, r *http.Request) {
	// Extract barcode dari URL: /api/produk/barcode/1234567890
	barcode := strings.TrimPrefix(r.URL.Path, "/api/produk/barcode/")
	if barcode == "" {
		log.Printf("⚠️ Handler: Empty barcode")
		http.Error(w, "Barcode tidak boleh kosong", http.StatusBadRequest)
		return
	}

	// Panggil service untuk cari produk by barcode
	product, err := h.service.GetByBarcode(barcode)
	if err != nil {
		log.Printf("❌ Handler: Error getting product by barcode %s: %v", barcode, err)
		http.Error(w, "Produk dengan barcode tersebut tidak ditemukan", http.StatusNotFound)
		return
	}

	// Filter sensitive data based on role
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		product.HargaBeli = nil
		product.Margin = nil
		product.CreatedBy = nil
	} else {
		product.Margin = product.CalculateMargin()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// Create adds a new product
// Fungsi ini handle POST /api/produk
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Check authorization
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, "Forbidden: Only Admin can create products", http.StatusForbidden)
		return
	}

	var product models.Product // Buat variable untuk menampung data dari request

	// Decode JSON dari request body ke struct product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		// Log error untuk debugging
		log.Printf("⚠️ Handler: Invalid request body for create product: %v", err)
		// Kalau JSON invalid, return error 400 (Bad Request)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set created_by from user
	userID := user.ID
	product.CreatedBy = &userID

	// Panggil service untuk create produk baru
	err = h.service.Create(&product)
	if err != nil {
		// Log sudah dilakukan di service layer
		// Kalau error validasi, return 400, kalau error lain return 500
		if strings.Contains(err.Error(), "tidak boleh") || strings.Contains(err.Error(), "harus") || strings.Contains(err.Error(), "minimal") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Set header Content-Type jadi application/json
	w.Header().Set("Content-Type", "application/json")

	// Set status code 201 (Created)
	w.WriteHeader(http.StatusCreated)

	// Encode product yang baru dibuat (sudah ada ID) dan kirim ke client
	json.NewEncoder(w).Encode(product)
}

// Update updates an existing product
// Fungsi ini handle PUT /api/produk/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Check authorization
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, "Forbidden: Only Admin can update products", http.StatusForbidden)
		return
	}

	// Extract ID dari URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// Convert string ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Log error untuk debugging
		log.Printf("⚠️ Handler: Invalid product ID for update: %s", idStr)
		// Kalau invalid ID, return error 400
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	var product models.Product // Buat variable untuk menampung data update

	// Decode JSON dari request body
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		// Log error untuk debugging
		log.Printf("⚠️ Handler: Invalid request body for update product ID %d: %v", id, err)
		// Kalau JSON invalid, return error 400
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Panggil service untuk update produk
	err = h.service.Update(id, &product)
	if err != nil {
		// Log sudah dilakukan di service layer
		// Kalau error validasi, return 400, kalau error lain return 500
		if strings.Contains(err.Error(), "tidak boleh") || strings.Contains(err.Error(), "harus") || strings.Contains(err.Error(), "minimal") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Set header Content-Type jadi application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode product yang sudah di-update dan kirim ke client
	json.NewEncoder(w).Encode(product)
}

// Delete removes a product
// Fungsi ini handle DELETE /api/produk/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Check authorization
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, "Forbidden: Only Admin can delete products", http.StatusForbidden)
		return
	}

	// Extract ID dari URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// Convert string ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Log error untuk debugging
		log.Printf("⚠️ Handler: Invalid product ID for delete: %s", idStr)
		// Kalau invalid ID, return error 400
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	// Panggil service untuk delete produk
	err = h.service.Delete(id)
	if err != nil {
		// Log sudah dilakukan di service layer
		// Kalau error saat delete, return error 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set header Content-Type jadi application/json
	w.Header().Set("Content-Type", "application/json")

	// Kirim response sukses delete
	json.NewEncoder(w).Encode(map[string]string{"message": "sukses delete"})
}

package handlers

import (
	"encoding/json"      // Package untuk encode/decode JSON
	"kasir-api/models"   // Import models untuk struct Category
	"kasir-api/services" // Import services untuk business logic
	"log"
	"net/http" // Package untuk HTTP server
	"strconv"  // Package untuk convert string ke int
	"strings"  // Package untuk manipulasi string
)

// CategoryHandler handles HTTP requests for categories
// Handler adalah layer yang berhadapan langsung dengan HTTP request/response
type CategoryHandler struct {
	service *services.CategoryService // Pointer ke CategoryService
}

// NewCategoryHandler creates a new CategoryHandler
// Fungsi ini adalah "constructor" untuk membuat instance CategoryHandler
func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service} // Return struct dengan service yang sudah di-inject
}

// HandleCategories handles /api/categories (GET all, POST new)
// Fungsi ini handle 2 method: GET (ambil semua) dan POST (tambah baru)
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	// Switch berdasarkan HTTP method
	switch r.Method {
	case "GET":
		h.GetAll(w, r) // Kalau GET, panggil GetAll
	case "POST":
		h.Create(w, r) // Kalau POST, panggil Create
	default:
		// Kalau method lain, return error
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleCategoryByID handles /api/categories/{id} (GET, PUT, DELETE)
// Fungsi ini handle 3 method: GET (by ID), PUT (update), DELETE (hapus)
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
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

// GetAll retrieves all categories
// Fungsi ini handle GET /api/categories
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Panggil service untuk ambil semua kategori
	categories, err := h.service.GetAll()
	if err != nil {
		// Log error untuk debugging
		log.Printf("❌ Handler: Error getting categories: %v", err)
		// Kalau error, kirim HTTP error 500 (Internal Server Error)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set header Content-Type jadi application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode categories jadi JSON dan kirim ke client
	json.NewEncoder(w).Encode(categories)
}

// GetByID retrieves a category by ID
// Fungsi ini handle GET /api/categories/{id}
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID dari URL
	// Misal URL: /api/categories/2 -> idStr = "2"
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// Convert string "2" jadi integer 2
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Log error untuk debugging
		log.Printf("⚠️ Handler: Invalid category ID: %s", idStr)
		// Kalau gagal convert, return error 400 (Bad Request)
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// Panggil service untuk ambil kategori by ID
	category, err := h.service.GetByID(id)
	if err != nil {
		// Log error untuk debugging
		log.Printf("❌ Handler: Error getting category ID %d: %v", id, err)
		// Kalau tidak ketemu, return error 404 (Not Found)
		http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
		return
	}

	// Set header Content-Type jadi application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode category jadi JSON dan kirim ke client
	json.NewEncoder(w).Encode(category)
}

// Create adds a new category
// Fungsi ini handle POST /api/categories
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category models.Category // Buat variable untuk menampung data dari request

	// Decode JSON dari request body ke struct category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		// Log error untuk debugging
		log.Printf("⚠️ Handler: Invalid request body for create category: %v", err)
		// Kalau JSON invalid, return error 400 (Bad Request)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Panggil service untuk create kategori baru
	err = h.service.Create(&category)
	if err != nil {
		// Log sudah dilakukan di service layer
		// Kalau error validasi, return 400, kalau error lain return 500
		if strings.Contains(err.Error(), "tidak boleh") || strings.Contains(err.Error(), "minimal") {
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

	// Encode category yang baru dibuat (sudah ada ID) dan kirim ke client
	json.NewEncoder(w).Encode(category)
}

// Update updates an existing category
// Fungsi ini handle PUT /api/categories/{id}
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extract ID dari URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// Convert string ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Log error untuk debugging
		log.Printf("⚠️ Handler: Invalid category ID for update: %s", idStr)
		// Kalau invalid ID, return error 400
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	var category models.Category // Buat variable untuk menampung data update

	// Decode JSON dari request body
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		// Log error untuk debugging
		log.Printf("⚠️ Handler: Invalid request body for update category ID %d: %v", id, err)
		// Kalau JSON invalid, return error 400
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Panggil service untuk update kategori
	err = h.service.Update(id, &category)
	if err != nil {
		// Log sudah dilakukan di service layer
		// Kalau error validasi, return 400, kalau error lain return 500
		if strings.Contains(err.Error(), "tidak boleh") || strings.Contains(err.Error(), "minimal") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Set header Content-Type jadi application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode category yang sudah di-update dan kirim ke client
	json.NewEncoder(w).Encode(category)
}

// Delete removes a category
// Fungsi ini handle DELETE /api/categories/{id}
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extract ID dari URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// Convert string ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Log error untuk debugging
		log.Printf("⚠️ Handler: Invalid category ID for delete: %s", idStr)
		// Kalau invalid ID, return error 400
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// Panggil service untuk delete kategori
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

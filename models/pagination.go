package models

// PaginationParams holds pagination parameters from query string
type PaginationParams struct {
	Page  int // Halaman ke berapa (mulai dari 1)
	Limit int // Berapa item per halaman
}

// PaginationMeta holds pagination metadata for response
type PaginationMeta struct {
	Page       int `json:"page"`        // Halaman saat ini
	Limit      int `json:"limit"`       // Limit per halaman
	TotalItems int `json:"total_items"` // Total semua items
	TotalPages int `json:"total_pages"` // Total halaman
}

// PaginatedResponse is a generic response with pagination
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`       // Data items (bisa []Product, []Category, dll)
	Pagination PaginationMeta `json:"pagination"` // Metadata pagination
}

// NewPaginationParams creates pagination params with defaults
func NewPaginationParams(page, limit int) PaginationParams {
	// Default page = 1
	if page < 1 {
		page = 1
	}

	// Default limit = 10, max = 100
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// GetOffset calculates the offset for SQL LIMIT/OFFSET
func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// CalculateTotalPages calculates total pages from total items
func CalculateTotalPages(totalItems, limit int) int {
	if limit == 0 {
		return 0
	}
	totalPages := totalItems / limit
	if totalItems%limit > 0 {
		totalPages++
	}
	return totalPages
}

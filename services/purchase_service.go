package services

import (
	"fmt"
	"kasir-api/models"
	"kasir-api/repositories"
	"log"
	"strings"
	"time"
)

// PurchaseService handles business logic for purchases
// Service layer untuk pembelian/pengadaan barang
type PurchaseService struct {
	repo         *repositories.PurchaseRepository
	productCache *CacheService // Untuk invalidate cache produk setelah pembelian
}

// NewPurchaseService creates a new PurchaseService
func NewPurchaseService(repo *repositories.PurchaseRepository, cache *CacheService) *PurchaseService {
	return &PurchaseService{
		repo:         repo,
		productCache: cache,
	}
}

// Create validates and creates a new purchase
// Fungsi ini memvalidasi data pembelian lalu menyimpannya
func (s *PurchaseService) Create(req *models.PurchaseRequest, createdBy int) (*models.Purchase, error) {
	// ─── VALIDASI ───

	// 1. Harus ada minimal 1 item
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("pembelian harus memiliki minimal 1 item")
	}

	// 2. Validasi setiap item
	for i, item := range req.Items {
		// Quantity harus > 0
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("item #%d: quantity harus lebih dari 0", i+1)
		}

		// BuyPrice harus >= 0
		if item.BuyPrice < 0 {
			return nil, fmt.Errorf("item #%d: harga beli tidak boleh negatif", i+1)
		}

		// Jika produk baru (tidak ada product_id)
		if item.ProductID == nil {
			// Nama produk wajib
			if item.ProductName == nil || strings.TrimSpace(*item.ProductName) == "" {
				return nil, fmt.Errorf("item #%d: nama produk wajib diisi untuk produk baru", i+1)
			}
			// Harga jual wajib
			if item.SellPrice == nil || *item.SellPrice <= 0 {
				return nil, fmt.Errorf("item #%d: harga jual wajib dan harus > 0 untuk produk baru", i+1)
			}
			// Harga jual harus > harga beli (peringatan jika rugi)
			if item.SellPrice != nil && *item.SellPrice <= item.BuyPrice {
				log.Printf("⚠️ Peringatan item #%d: harga jual (%.0f) <= harga beli (%.0f), margin negatif!",
					i+1, *item.SellPrice, item.BuyPrice)
			}
		}
	}

	// ─── SIMPAN KE DATABASE ───
	purchase, err := s.repo.Create(req, createdBy)
	if err != nil {
		log.Printf("❌ Error creating purchase: %v", err)
		return nil, err
	}

	// ─── INVALIDATE CACHE ───
	// Karena stok dan harga_beli produk berubah, cache produk harus di-clear
	s.productCache.DeletePattern("products:*")
	log.Printf("🗑️ Cache produk di-invalidate setelah pembelian")

	return purchase, nil
}

// GetAll retrieves all purchases
// Fungsi ini mengambil riwayat semua pembelian
func (s *PurchaseService) GetAll() ([]models.Purchase, error) {
	purchases, err := s.repo.GetAll()
	if err != nil {
		log.Printf("❌ Error getting purchases: %v", err)
		return nil, err
	}
	return purchases, nil
}

// GetByID retrieves a purchase by ID with items
// Fungsi ini mengambil detail 1 pembelian
func (s *PurchaseService) GetByID(id int) (*models.Purchase, error) {
	purchase, err := s.repo.GetByID(id)
	if err != nil {
		log.Printf("❌ Error getting purchase ID %d: %v", id, err)
		return nil, err
	}
	return purchase, nil
}

// GetTotalPengeluaran retrieves total expenditure for a date range
// Fungsi ini menghitung total pengeluaran untuk laporan
func (s *PurchaseService) GetTotalPengeluaran(startDate, endDate time.Time) (float64, int, error) {
	return s.repo.GetTotalPengeluaran(startDate, endDate)
}

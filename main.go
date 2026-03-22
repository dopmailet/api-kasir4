package main

import (
	"encoding/json"          // Package untuk encode/decode JSON
	"fmt"                    // Package untuk print ke console
	"kasir-api/config"       // Import package config untuk configuration management
	"kasir-api/database"     // Import package database untuk koneksi DB
	"kasir-api/handlers"     // Import package handlers untuk HTTP handlers
	"kasir-api/middleware"   // Import package middleware untuk auth, logging, CORS
	"kasir-api/repositories" // Import package repositories untuk database operations
	"kasir-api/services"     // Import package services untuk business logic
	"log"                    // Package untuk logging
	"log/slog"               // Package untuk structured logging
	"net/http"               // Package untuk HTTP server
	"os"                     // Package untuk environment variables
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	// ==================== SETUP LOGGING ====================
	// Setup structured logging dengan slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// ==================== LOAD CONFIGURATION ====================
	// Load config dari .env file dan environment variables menggunakan Viper
	// Setup env explicitly using godotenv
	godotenv.Load(".env")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("❌ Failed to load configuration:", err)
	}

	// Verify JWT_SECRET is set
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("❌ JWT_SECRET is not set in environment variables!")
	}

	// ==================== INITIALIZE DATABASE ====================
	// Panggil InitDB untuk connect ke database PostgreSQL
	// Gunakan connection string dari config
	db := database.InitDB(cfg.GetDatabaseURL())

	// defer = pastikan db.Close() dipanggil saat program selesai
	// Ini penting untuk tutup koneksi database dengan benar
	defer db.Close()

	// ==================== INITIALIZE REDIS ====================
	// Initialize Redis connection untuk caching
	config.InitRedis()
	// defer = pastikan Redis connection ditutup saat program selesai
	defer config.CloseRedis()

	// ==================== DEPENDENCY INJECTION ====================
	// Dependency Injection = "inject" dependency ke layer yang membutuhkan
	// Flow: Database -> Repository -> Service -> Handler

	// Cache service (shared across all services)
	cacheService := services.NewCacheService()

	// User & Auth layers (NEW!)
	storeRepo := repositories.NewStoreRepository(db)
	pkgRepo := repositories.NewSubscriptionPackageRepository(db)
	storeService := services.NewStoreService(storeRepo, pkgRepo)

	userRepo := repositories.NewUserRepository(db)                  // Inject db ke repository
	authService := services.NewAuthService(db, userRepo, storeRepo) // Inject db, userRepo, storeRepo ke service
	authHandler := handlers.NewAuthHandler(authService)             // Inject service ke handler
	userService := services.NewUserService(userRepo, storeService)  // Inject repo ke service
	userHandler := handlers.NewUserHandler(userService)             // Inject service ke handler

	// Superadmin layers (Platform Manager Only)
	superadminService := services.NewSuperadminService(storeRepo, pkgRepo, userRepo)
	superadminHandler := handlers.NewSuperadminHandler(superadminService)

	// Product layers
	productRepo := repositories.NewProductRepository(db)                                  // Inject db ke repository
	productService := services.NewProductService(productRepo, cacheService, storeService) // Inject repo, cache, and storeSvc ke service
	productHandler := handlers.NewProductHandler(productService)                          // Inject service ke handler

	// Category layers
	categoryRepo := repositories.NewCategoryRepository(db)                     // Inject db ke repository
	categoryService := services.NewCategoryService(categoryRepo, cacheService) // Inject repo dan cache ke service
	categoryHandler := handlers.NewCategoryHandler(categoryService)            // Inject service ke handler

	// Transaction layers
	transactionRepo := repositories.NewTransactionRepository(db)                        // Inject db ke repository
	transactionService := services.NewTransactionService(transactionRepo, storeService) // Inject repo dan storeSvc ke service
	transactionHandler := handlers.NewTransactionHandler(transactionService)            // Inject service ke handler

	// Report layers
	reportRepo := repositories.NewReportRepository(db)        // Inject db ke repository
	reportService := services.NewReportService(reportRepo)    // Inject repo ke service
	reportHandler := handlers.NewReportHandler(reportService) // Inject service ke handler

	// Discount layers
	discountRepo := repositories.NewDiscountRepository(db)
	discountHandler := handlers.NewDiscountHandler(discountRepo)

	// Purchase layers (Admin Only)
	purchaseRepo := repositories.NewPurchaseRepository(db)
	purchaseService := services.NewPurchaseService(purchaseRepo, cacheService, storeService)
	purchaseHandler := handlers.NewPurchaseHandler(purchaseService)

	// Employee layers (Admin Only)
	employeeRepo := repositories.NewEmployeeRepository(db)
	employeeService := services.NewEmployeeService(employeeRepo)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

	// Payroll layers (Admin Only)
	payrollRepo := repositories.NewPayrollRepository(db)
	payrollService := services.NewPayrollService(payrollRepo)
	payrollHandler := handlers.NewPayrollHandler(payrollService)

	// Expense layers (Admin Only)
	expenseRepo := repositories.NewExpenseRepository(db)
	expenseService := services.NewExpenseService(expenseRepo)
	expenseHandler := handlers.NewExpenseHandler(expenseService)

	// Cash Flow layers (Admin Only)
	cashFlowRepo := repositories.NewCashFlowRepository(db)
	cashFlowService := services.NewCashFlowService(cashFlowRepo)
	cashFlowHandler := handlers.NewCashFlowHandler(cashFlowService)

	// ==================== VALIDATOR ====================
	validate := validator.New()

	// Customer & Loyalty layers
	customerRepo := repositories.NewCustomerRepository(db)
	loyaltyRepo := repositories.NewLoyaltyRepository(db)
	customerService := services.NewCustomerService(customerRepo, loyaltyRepo)
	customerHandler := handlers.NewCustomerHandler(customerService, validate)

	// Settings layers
	settingRepo := repositories.NewSettingRepository(db)
	settingService := services.NewSettingService(settingRepo)
	settingHandler := handlers.NewSettingHandler(settingService)

	// Supplier & Payables layers (Admin Only)
	supplierRepo := repositories.NewSupplierRepository(db)
	supplierService := services.NewSupplierService(supplierRepo)
	supplierHandler := handlers.NewSupplierHandler(supplierService)

	// Store layers (Tenant Info)
	storeHandler := handlers.NewStoreHandler(storeService)

	// ==================== SETUP ROUTER WITH MIDDLEWARE ====================
	// Create a new ServeMux for better routing
	mux := http.NewServeMux()

	// Employee routes (Admin Only)
	mux.Handle("/api/employees", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			employeeHandler.GetAll(w, r)
		} else if r.Method == http.MethodPost {
			employeeHandler.Create(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	mux.Handle("/api/employees/", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			employeeHandler.GetByID(w, r)
		} else if r.Method == http.MethodPut {
			employeeHandler.Update(w, r)
		} else if r.Method == http.MethodDelete {
			employeeHandler.SoftDelete(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	// Payroll routes (Admin Only)
	mux.Handle("/api/payroll/report", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			payrollHandler.GetReport(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	mux.Handle("/api/payroll", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			payrollHandler.GetAll(w, r)
		} else if r.Method == http.MethodPost {
			payrollHandler.Create(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	mux.Handle("/api/payroll/", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			payrollHandler.GetByID(w, r)
		} else if r.Method == http.MethodPut {
			payrollHandler.Update(w, r)
		} else if r.Method == http.MethodDelete {
			payrollHandler.Delete(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	// Expense routes (Admin Only)
	mux.Handle("/api/expenses", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			expenseHandler.GetAll(w, r)
		} else if r.Method == http.MethodPost {
			expenseHandler.Create(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	mux.Handle("/api/expenses/", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			expenseHandler.GetByID(w, r)
		} else if r.Method == http.MethodPut {
			expenseHandler.Update(w, r)
		} else if r.Method == http.MethodDelete {
			expenseHandler.Delete(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	// Purchase routes (Admin Only)
	// /api/purchases -> GET (list), POST (create)
	mux.Handle("/api/purchases/", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(purchaseHandler.HandlePurchaseByID))))
	mux.Handle("/api/purchases", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(purchaseHandler.HandlePurchases))))

	// Dashboard routes
	// /api/dashboard/sales-trend -> GET (Admin Only) ?period=day|month|year&start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
	mux.Handle("/api/dashboard/sales-trend", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(reportHandler.GetSalesTrend))))
	// /api/dashboard/top-products -> GET (Admin Only) ?limit=5
	mux.Handle("/api/dashboard/top-products", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(reportHandler.GetTopProducts))))
	// /api/dashboard/summary -> GET (Admin Only) ?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&low_stock_threshold=5
	mux.Handle("/api/dashboard/summary", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(reportHandler.GetDashboardSummary))))
	// /api/dashboard/assets -> GET (Admin Only)
	mux.Handle("/api/dashboard/assets", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(reportHandler.GetDashboardAssets))))

	// Cash Flow routes
	// /api/cash-flow/summary -> GET (Admin Only) ?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&timezone=Asia/Jakarta
	mux.Handle("/api/cash-flow/summary", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(cashFlowHandler.GetSummary))))
	// /api/cash-flow/trend -> GET (Admin Only) ?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&timezone=Asia/Jakarta
	mux.Handle("/api/cash-flow/trend", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(cashFlowHandler.GetTrend))))

	// Discount routes
	// /api/discounts/active -> GET (Public/Kasir)
	mux.Handle("/api/discounts/active", middleware.AuthMiddleware(http.HandlerFunc(discountHandler.GetActive)))

	// /api/discounts -> GET (Admin), POST (Admin)
	mux.Handle("/api/discounts", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			discountHandler.GetAll(w, r)
		} else if r.Method == http.MethodPost {
			discountHandler.Create(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	// /api/discounts/ -> PUT, DELETE (Admin)
	mux.Handle("/api/discounts/", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handler ini akan menangkap /api/discounts/{id}
		// Cek method
		if r.Method == http.MethodPut {
			discountHandler.Update(w, r)
		} else if r.Method == http.MethodDelete {
			discountHandler.Delete(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	// Settings routes
	// /api/settings/customer -> GET, PUT (Admin)
	mux.Handle("/api/settings/customer", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			settingHandler.GetCustomerSettings(w, r)
		} else if r.Method == http.MethodPut {
			settingHandler.UpdateCustomerSettings(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	// Customer routes
	// /api/customers/search -> GET (Admin, Kasir)
	mux.Handle("/api/customers/search", middleware.AuthMiddleware(http.HandlerFunc(customerHandler.Search)))

	// /api/customers -> GET (Admin for list), POST (Admin, Kasir)
	mux.Handle("/api/customers", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Require Admin for full list
			middleware.RequireAdmin(http.HandlerFunc(customerHandler.GetAll)).ServeHTTP(w, r)
		} else if r.Method == http.MethodPost {
			// Admin & Kasir can create
			customerHandler.Create(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// /api/customers/{id} -> GET, PUT (Admin for PUT, Admin/Kasir for GET)
	mux.Handle("/api/customers/", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// Check for specific sub-routes like /transactions and /loyalty-transactions
		if strings.HasSuffix(path, "/transactions") {
			customerHandler.GetTransactions(w, r)
			return
		} else if strings.HasSuffix(path, "/loyalty-transactions") {
			customerHandler.GetLoyaltyHistory(w, r)
			return
		}

		// Handle base /api/customers/{id}
		if r.Method == http.MethodGet {
			customerHandler.GetByID(w, r)
		} else if r.Method == http.MethodPut {
			middleware.RequireAdmin(http.HandlerFunc(customerHandler.Update)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// ==================== PUBLIC ROUTES (No Auth Required) ====================

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Kasir API Running - Session 4 with Authentication",
			"version": "1.0.0",
		})
	})

	// Auth routes (public - no auth required)
	// /api/auth/login
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodOptions {
			authHandler.Login(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// /api/auth/register (public)
	mux.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodOptions {
			authHandler.Register(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Paket Langganan (public - untuk landing page)
	mux.HandleFunc("/api/packages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodOptions {
			superadminHandler.GetPublicPackages(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Pengaturan Platform (public - untuk widget badge langganan via WA)
	mux.HandleFunc("/api/settings/public", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodOptions {
			settingHandler.GetPublicSettings(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// ==================== PROTECTED ROUTES (Auth Required) ====================

	// Middleware for authentication
	// We wrap the handlers with AuthMiddleware to ensure only authenticated users can access

	// User routes (Admin Only)
	mux.Handle("/api/users", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.GetAll(w, r)
		} else if r.Method == http.MethodPost {
			userHandler.Create(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	mux.Handle("/api/users/", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut && strings.HasSuffix(r.URL.Path, "/status") {
			userHandler.UpdateStatus(w, r)
		} else if r.Method == http.MethodPut && strings.HasSuffix(r.URL.Path, "/password") {
			userHandler.UpdatePassword(w, r)
		} else if r.Method == http.MethodDelete {
			userHandler.Delete(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	// Product routes
	// /api/produk/ -> GET (by ID), PUT, DELETE
	mux.Handle("/api/produk/", middleware.AuthMiddleware(http.HandlerFunc(productHandler.HandleProductByID)))

	// /api/produk -> GET (all), POST
	mux.Handle("/api/produk", middleware.AuthMiddleware(http.HandlerFunc(productHandler.HandleProducts)))

	// Category routes
	mux.Handle("/api/categories/", middleware.AuthMiddleware(http.HandlerFunc(categoryHandler.HandleCategoryByID)))
	mux.Handle("/api/categories", middleware.AuthMiddleware(http.HandlerFunc(categoryHandler.HandleCategories)))

	// Transaction routes
	mux.Handle("/api/checkout", middleware.AuthMiddleware(http.HandlerFunc(transactionHandler.Checkout)))
	mux.Handle("/api/transactions/", middleware.AuthMiddleware(http.HandlerFunc(transactionHandler.HandleTransactionByID))) // GET by ID
	mux.Handle("/api/transactions", middleware.AuthMiddleware(http.HandlerFunc(transactionHandler.HandleTransactions)))     // GET all

	// Report routes
	mux.Handle("/api/report/hari-ini", middleware.AuthMiddleware(http.HandlerFunc(reportHandler.GetDailySalesReport)))
	mux.Handle("/api/report", middleware.AuthMiddleware(http.HandlerFunc(reportHandler.GetSalesReportByDateRange)))

	// Cash Flow ledger route (Admin Only) — summary & trend sudah didaftarkan di atas (line 238-240)
	mux.Handle("/api/cash-flow/ledger", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(cashFlowHandler.GetLedger))))

	// Cash Flow funds routes (Dana & Modal)
	mux.Handle("/api/cash-flow/funds", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cashFlowHandler.GetFunds(w, r)
		case http.MethodPost:
			cashFlowHandler.CreateFund(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	mux.Handle("/api/cash-flow/funds/", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			cashFlowHandler.DeleteFund(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	mux.Handle("/api/cash-flow/initial-balance", middleware.AuthMiddleware(http.HandlerFunc(cashFlowHandler.GetInitialBalance)))

	// Superadmin routes (Platform Manager Only)
	mux.Handle("/api/superadmin/stores", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			superadminHandler.GetAllStores(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/superadmin/stores/bulk", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			superadminHandler.DeleteBulkStores(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/superadmin/stores/", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut && strings.HasSuffix(r.URL.Path, "/status") {
			superadminHandler.UpdateStoreStatus(w, r)
		} else if r.Method == http.MethodPut && strings.HasSuffix(r.URL.Path, "/package") {
			superadminHandler.UpgradePackage(w, r)
		} else if r.Method == http.MethodDelete {
			superadminHandler.DeleteStore(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/superadmin/packages", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			superadminHandler.GetAllPackages(w, r)
		} else if r.Method == http.MethodPost {
			superadminHandler.CreatePackage(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/superadmin/packages/", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			superadminHandler.UpdatePackage(w, r)
		} else if r.Method == http.MethodPatch && strings.HasSuffix(r.URL.Path, "/popular") {
			superadminHandler.TogglePackagePopular(w, r)
		} else if r.Method == http.MethodDelete {
			superadminHandler.DeletePackage(w, r)
		} else if r.Method == http.MethodGet {
			superadminHandler.GetPackageByID(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/superadmin/settings", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			settingHandler.GetPlatformSettings(w, r)
		} else if r.Method == http.MethodPut {
			settingHandler.UpdatePlatformSettings(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/superadmin/store-users", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			superadminHandler.GetStoreUsers(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/store/info", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			storeHandler.GetMyStoreInfo(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/store/limits", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			storeHandler.GetLimits(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// ==================== APPLY GLOBAL MIDDLEWARE ====================
	// Middleware chain: CORS -> Logging -> Handler
	handler := middleware.CORSMiddleware(
		middleware.LoggingMiddleware(mux),
	)

	// ==================== START SERVER ====================
	port := cfg.Port

	// Print informasi server
	fmt.Println("🚀 ========================================")
	fmt.Println("🚀 Kasir API - Session 4 (Authentication)")
	fmt.Println("🚀 ========================================")
	fmt.Println("📡 Server running on port:", port)
	fmt.Println("🔐 Authentication: ENABLED")
	fmt.Println("📝 Logging: ENABLED (structured JSON)")
	fmt.Println("🌐 CORS: ENABLED")
	fmt.Println("")
	fmt.Println("📚 Employee & Payroll Endpoints (Admin Only):")
	fmt.Println("  - GET    /api/employees")
	fmt.Println("  - POST   /api/employees")
	fmt.Println("  - GET    /api/employees/{id}")
	fmt.Println("  - PUT    /api/employees/{id}")
	fmt.Println("  - DELETE /api/employees/{id}")
	fmt.Println("  - GET    /api/payroll")
	fmt.Println("  - POST   /api/payroll")
	fmt.Println("  - GET    /api/payroll/report")
	fmt.Println("  - GET    /api/payroll/{id}")
	fmt.Println("  - PUT    /api/payroll/{id}")
	fmt.Println("  - DELETE /api/payroll/{id}")
	fmt.Println("")
	fmt.Println("📚 Public Endpoints:")
	fmt.Println("  - GET    /health")
	fmt.Println("  - POST   /api/auth/login")
	fmt.Println("  - POST   /api/auth/register")
	fmt.Println("")
	fmt.Println("📚 User Endpoints (Admin Only):")
	fmt.Println("  - GET    /api/users")
	fmt.Println("  - POST   /api/users")
	fmt.Println("  - PUT    /api/users/{id}/password")
	fmt.Println("  - DELETE /api/users/{id}")
	fmt.Println("")
	fmt.Println("📚 Superadmin Endpoints:")
	fmt.Println("  - GET    /api/superadmin/stores")
	fmt.Println("  - PUT    /api/superadmin/stores/{id}/status")
	fmt.Println("  - GET    /api/superadmin/packages")
	fmt.Println("")
	fmt.Println("📚 Product Endpoints:")
	fmt.Println("  - GET    /api/produk")
	fmt.Println("  - GET    /api/produk?barcode=xxx")
	fmt.Println("  - POST   /api/produk")
	fmt.Println("  - GET    /api/produk/{id}")
	fmt.Println("  - GET    /api/produk/barcode/{code}")
	fmt.Println("  - PUT    /api/produk/{id}")
	fmt.Println("  - DELETE /api/produk/{id}")
	fmt.Println("")
	fmt.Println("📚 Category Endpoints:")
	fmt.Println("  - GET    /api/categories")
	fmt.Println("  - POST   /api/categories")
	fmt.Println("  - GET    /api/categories/{id}")
	fmt.Println("  - PUT    /api/categories/{id}")
	fmt.Println("  - DELETE /api/categories/{id}")
	fmt.Println("")
	fmt.Println("📚 Transaction Endpoints:")
	fmt.Println("  - POST   /api/checkout")
	fmt.Println("  - GET    /api/transactions")
	fmt.Println("")
	fmt.Println("📚 Purchase Endpoints (Admin Only):")
	fmt.Println("  - POST   /api/purchases")
	fmt.Println("  - GET    /api/purchases")
	fmt.Println("  - GET    /api/purchases/{id}")
	fmt.Println("")
	fmt.Println("📚 Report & Dashboard Endpoints (Admin Only):")
	fmt.Println("  - GET    /api/report/hari-ini")
	fmt.Println("  - GET    /api/report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD")
	fmt.Println("  - GET    /api/dashboard/summary?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&low_stock_threshold=5")
	fmt.Println("  - GET    /api/dashboard/sales-trend?period=day|month|year&start_date=YYYY-MM-DD&end_date=YYYY-MM-DD")
	fmt.Println("  - GET    /api/dashboard/top-products?limit=5")
	fmt.Println("")
	fmt.Println("🔑 Default Credentials:")
	fmt.Println("  - admin / admin123 (role: admin)")
	fmt.Println("  - kasir1 / kasir123 (role: kasir)")
	fmt.Println("")
	// ─── Supplier Routes (Admin Only) ───
	// PENTING: /api/suppliers/debt-summary harus SEBELUM /api/suppliers/ agar tidak ditangkap wildcard
	mux.Handle("/api/suppliers/debt-summary", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			supplierHandler.GetDebtSummary(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	mux.Handle("/api/suppliers", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			supplierHandler.GetAll(w, r)
		case http.MethodPost:
			supplierHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	mux.Handle("/api/suppliers/", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// /api/suppliers/{id}/payables  or  /api/suppliers/{id}/payables/{payableId}
		if strings.Contains(path, "/payables/") {
			// PUT /api/suppliers/{id}/payables/{payableId}
			if r.Method == http.MethodPut {
				supplierHandler.UpdatePayable(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else if strings.HasSuffix(path, "/payables") {
			// GET or POST /api/suppliers/{id}/payables
			switch r.Method {
			case http.MethodGet:
				supplierHandler.GetPayables(w, r)
			case http.MethodPost:
				supplierHandler.CreatePayable(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			// /api/suppliers/{id}
			switch r.Method {
			case http.MethodGet:
				supplierHandler.GetByID(w, r)
			case http.MethodPut:
				supplierHandler.Update(w, r)
			case http.MethodDelete:
				supplierHandler.Delete(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	}))))

	// ─── Payable Payment Routes (Admin Only) ───
	mux.Handle("/api/payables/", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			supplierHandler.GetPayments(w, r)
		case http.MethodPost:
			supplierHandler.CreatePayment(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	// /api/payable-payments -> GET (Admin Only) ?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
	mux.Handle("/api/payable-payments", middleware.AuthMiddleware(middleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			supplierHandler.GetAllPayablePayments(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	fmt.Println("✅ Ready to accept requests!")
	fmt.Println("========================================")

	// Start HTTP server
	slog.Info("Starting server", "port", port)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		fmt.Println("❌ Failed to start server:", err)
	}
}

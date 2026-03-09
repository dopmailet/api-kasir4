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
)

func main() {
	// ==================== SETUP LOGGING ====================
	// Setup structured logging dengan slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// ==================== LOAD CONFIGURATION ====================
	// Load config dari .env file dan environment variables menggunakan Viper
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
	userRepo := repositories.NewUserRepository(db)      // Inject db ke repository
	authService := services.NewAuthService(userRepo)    // Inject repo ke service
	authHandler := handlers.NewAuthHandler(authService) // Inject service ke handler
	userService := services.NewUserService(userRepo)    // Inject repo ke service
	userHandler := handlers.NewUserHandler(userService) // Inject service ke handler

	// Product layers
	productRepo := repositories.NewProductRepository(db)                    // Inject db ke repository
	productService := services.NewProductService(productRepo, cacheService) // Inject repo dan cache ke service
	productHandler := handlers.NewProductHandler(productService)            // Inject service ke handler

	// Category layers
	categoryRepo := repositories.NewCategoryRepository(db)                     // Inject db ke repository
	categoryService := services.NewCategoryService(categoryRepo, cacheService) // Inject repo dan cache ke service
	categoryHandler := handlers.NewCategoryHandler(categoryService)            // Inject service ke handler

	// Transaction layers
	transactionRepo := repositories.NewTransactionRepository(db)             // Inject db ke repository
	transactionService := services.NewTransactionService(transactionRepo)    // Inject repo ke service
	transactionHandler := handlers.NewTransactionHandler(transactionService) // Inject service ke handler

	// Report layers
	reportRepo := repositories.NewReportRepository(db)        // Inject db ke repository
	reportService := services.NewReportService(reportRepo)    // Inject repo ke service
	reportHandler := handlers.NewReportHandler(reportService) // Inject service ke handler

	// Discount layers
	discountRepo := repositories.NewDiscountRepository(db)
	discountHandler := handlers.NewDiscountHandler(discountRepo)

	// Purchase layers (Admin Only)
	purchaseRepo := repositories.NewPurchaseRepository(db)
	purchaseService := services.NewPurchaseService(purchaseRepo, cacheService)
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
	mux.HandleFunc("/api/auth/login", authHandler.Login)

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
		if r.Method == http.MethodPut && strings.HasSuffix(r.URL.Path, "/password") {
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
	fmt.Println("")
	fmt.Println("📚 User Endpoints (Admin Only):")
	fmt.Println("  - GET    /api/users")
	fmt.Println("  - POST   /api/users")
	fmt.Println("  - PUT    /api/users/{id}/password")
	fmt.Println("  - DELETE /api/users/{id}")
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

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found or error loading")
	}

	dbUrl := os.Getenv("DB_CONN")
	if dbUrl == "" {
		dbUrl = "postgresql://postgres.bwpaqpblhnsmrtuivvyf:1LDjhKBAmgvhK61s@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres"
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get Customer 1
	var loyaltyPoints, totalSpent, totalTransactions int
	var isActive bool
	err = db.QueryRow("SELECT loyalty_points, total_spent, total_transactions, is_active FROM customers WHERE id = 1").Scan(
		&loyaltyPoints, &totalSpent, &totalTransactions, &isActive,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Customer ID 1 in DB:\n  Loyalty Points: %d\n  Total Spent: %d\n  Total Transactions: %d\n  Is Active: %v\n",
		loyaltyPoints, totalSpent, totalTransactions, isActive)

	// Check loyalty transactions
	rows, err := db.Query("SELECT id, points, type, store_id FROM loyalty_transactions WHERE customer_id = 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("\nLoyalty Transactions for Customer 1:")
	for rows.Next() {
		var id, points, storeID int
		var lType string
		rows.Scan(&id, &points, &lType, &storeID)
		fmt.Printf("  ID: %d, Points: %d, Type: %s, StoreID: %d\n", id, points, lType, storeID)
	}
}

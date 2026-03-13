package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	timezone := "Asia/Makassar"

	query := `
		SELECT CURRENT_DATE AT TIME ZONE $1 as tz_date, 
		       (CURRENT_DATE AT TIME ZONE $1) + INTERVAL '1 day' as tz_next_date,
		       NOW() as now_raw,
		       NOW() AT TIME ZONE $1 as now_tz
	`

	var tzDate, tzNextDate, nowRaw, nowTz string
	err = db.QueryRow(query, timezone).Scan(&tzDate, &tzNextDate, &nowRaw, &nowTz)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	fmt.Printf("Timezone: %s\n", timezone)
	fmt.Printf("CURRENT_DATE AT TIME ZONE: %s\n", tzDate)
	fmt.Printf("CURRENT_DATE AT TIME ZONE + 1 day: %s\n", tzNextDate)
	fmt.Printf("NOW(): %s\n", nowRaw)
	fmt.Printf("NOW() AT TIME ZONE: %s\n", nowTz)

	// Test transactions query
	storeID := 1 // Assume store_id=1 for tested store
	var count int
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM transactions 
		WHERE store_id = $1 
		  AND created_at >= (CURRENT_DATE AT TIME ZONE $2)
		  AND created_at < (CURRENT_DATE AT TIME ZONE $2) + INTERVAL '1 day'
	`, storeID, timezone).Scan(&count)
	if err != nil {
		log.Fatalf("Count query failed: %v", err)
	}
	fmt.Printf("Transactions Today (store %d, %s) (CountTodayTransactions style): %d\n", storeID, timezone, count)

	// Compare with UTC boundary just to be sure
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM transactions 
		WHERE store_id = $1 
		  AND created_at >= CURRENT_DATE
		  AND created_at < CURRENT_DATE + INTERVAL '1 day'
	`, storeID).Scan(&count)
	if err != nil {
		log.Fatalf("Count query failed: %v", err)
	}
	fmt.Printf("Transactions Today (store %d, CURRENT_DATE literal): %d\n", storeID, count)

	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM transactions 
		WHERE store_id = $1 
		  AND DATE(created_at AT TIME ZONE 'UTC' AT TIME ZONE $2) = CURRENT_DATE
	`, storeID, timezone).Scan(&count)
	if err != nil {
		log.Fatalf("Count query failed: %v", err)
	}
	fmt.Printf("Transactions Today (store %d, DATE() literal): %d\n", storeID, count)
}

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := "postgresql://postgres.rgokcvgkeixqjisurvlh:ma2b4CiwqzjFn7iD@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("=== Tables in new database ===")
	count := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- %s\n", name)
		count++
	}
	fmt.Printf("Total tables: %d\n", count)
}

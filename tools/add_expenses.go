package main

import (
	"context"
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

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Connected to new Supabase database!")

	// The expenses table
	stmt := `CREATE TABLE expenses (
		id SERIAL PRIMARY KEY,
		category character varying(50) NOT NULL,
		description character varying(255) NOT NULL,
		amount numeric NOT NULL,
		expense_date date NOT NULL,
		is_recurring boolean DEFAULT false,
		recurring_period character varying(20),
		notes text,
		created_by integer,
		created_at timestamp without time zone DEFAULT now(),
		updated_at timestamp without time zone DEFAULT now()
	);`

	fmt.Println("Creating expenses table...")
	_, err = db.ExecContext(context.Background(), stmt)
	if err != nil {
		log.Printf("⚠️ Error creating expenses table: %v\n", err)
	} else {
		fmt.Printf("✅ Created expenses table successfully!\n")
	}
}

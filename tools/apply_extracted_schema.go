package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

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

	// 1. Wipe everything first
	fmt.Println("Wiping existing schema to avoid conflicts...")
	_, err = db.ExecContext(context.Background(), `
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
		GRANT ALL ON SCHEMA public TO postgres;
		GRANT ALL ON SCHEMA public TO public;
	`)
	if err != nil {
		log.Fatalf("Failed to wipe schema: %v", err)
	}

	// 2. Apply extracted schema
	content, err := ioutil.ReadFile("extracted_schema.sql")
	if err != nil {
		log.Fatalf("Failed to read extracted schema: %v", err)
	}

	fmt.Println("Applying schema to the new database...")

	// Split by double newline to run each CREATE TABLE separately
	statements := strings.Split(string(content), ";\n\n")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		_, err = db.ExecContext(context.Background(), stmt)
		if err != nil {
			log.Printf("⚠️ Error executing table creation:\n%s\nError: %v\n", stmt, err)
		} else {
			fmt.Printf("✅ Created table successfully!\n")
		}
	}

	fmt.Println("\n🎉 All tables successfully cloned from old database!")
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := "postgresql://postgres.bprzmdcmzqjlwiidmkfd:bfZmM64424rLdsDl@aws-1-ap-south-1.pooler.supabase.com:6543/postgres"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get all tables
	rows, err := db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		tables = append(tables, name)
	}

	outFile, err := os.Create("extracted_schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	outFile.WriteString("-- EXTRACTED SCHEMA FROM LIMA DATABASE --\n\n")

	for _, table := range tables {
		fmt.Printf("Extracting table: %s\n", table)

		// Get columns
		colsRows, err := db.QueryContext(context.Background(), `
			SELECT column_name, data_type, character_maximum_length, column_default, is_nullable
			FROM information_schema.columns
			WHERE table_schema = 'public' AND table_name = $1
			ORDER BY ordinal_position
		`, table)
		if err != nil {
			log.Fatal(err)
		}

		var columns []string
		for colsRows.Next() {
			var colName, dataType, isNullable string
			var charMaxLen, colDefault sql.NullString

			err := colsRows.Scan(&colName, &dataType, &charMaxLen, &colDefault, &isNullable)
			if err != nil {
				log.Fatal(err)
			}

			// Fix sequence defaults by just using SERIAL for ID columns
			if colName == "id" && strings.Contains(colDefault.String, "nextval") {
				colDef := fmt.Sprintf("    %s SERIAL PRIMARY KEY", colName)
				columns = append(columns, colDef)
				continue
			}

			colDef := fmt.Sprintf("    %s %s", colName, dataType)
			if charMaxLen.Valid {
				colDef += fmt.Sprintf("(%s)", charMaxLen.String)
			}
			if colDefault.Valid {
				colDef += fmt.Sprintf(" DEFAULT %s", colDefault.String)
			}
			if isNullable == "NO" {
				colDef += " NOT NULL"
			}
			columns = append(columns, colDef)
		}
		colsRows.Close()

		createTable := fmt.Sprintf("CREATE TABLE %s (\n%s\n);\n\n", table, strings.Join(columns, ",\n"))
		outFile.WriteString(createTable)
	}

	fmt.Println("Schema extracted to extracted_schema.sql")
}

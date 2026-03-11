package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Ping()

	// Dump semua tabel dan kolomnya
	query := `
		SELECT 
			t.table_name,
			c.column_name,
			c.data_type,
			c.character_maximum_length,
			c.column_default,
			c.is_nullable,
			c.udt_name
		FROM information_schema.tables t
		JOIN information_schema.columns c ON c.table_name = t.table_name AND c.table_schema = t.table_schema
		WHERE t.table_schema = 'public' AND t.table_type = 'BASE TABLE'
		ORDER BY t.table_name, c.ordinal_position
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	currentTable := ""
	for rows.Next() {
		var tableName, colName, dataType, isNullable, udtName string
		var maxLen sql.NullInt64
		var colDefault sql.NullString

		rows.Scan(&tableName, &colName, &dataType, &maxLen, &colDefault, &isNullable, &udtName)

		if tableName != currentTable {
			if currentTable != "" {
				fmt.Println()
			}
			fmt.Printf("=== TABLE: %s ===\n", tableName)
			currentTable = tableName
		}

		typeStr := udtName
		if maxLen.Valid {
			typeStr = fmt.Sprintf("%s(%d)", udtName, maxLen.Int64)
		}
		defStr := ""
		if colDefault.Valid {
			defStr = fmt.Sprintf(" DEFAULT %s", colDefault.String)
		}
		nullStr := ""
		if isNullable == "NO" {
			nullStr = " NOT NULL"
		}

		fmt.Printf("  %-25s %-20s%s%s\n", colName, typeStr, nullStr, defStr)
	}

	// Dump constraints (PK, FK, UNIQUE)
	fmt.Println("\n\n=== CONSTRAINTS ===")
	cquery := `
		SELECT 
			tc.table_name,
			tc.constraint_name,
			tc.constraint_type,
			kcu.column_name,
			ccu.table_name AS foreign_table,
			ccu.column_name AS foreign_column
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu ON kcu.constraint_name = tc.constraint_name AND kcu.table_schema = tc.table_schema
		LEFT JOIN information_schema.constraint_column_usage ccu ON ccu.constraint_name = tc.constraint_name AND ccu.table_schema = tc.table_schema
		WHERE tc.table_schema = 'public'
		ORDER BY tc.table_name, tc.constraint_type
	`
	rows2, _ := db.Query(cquery)
	defer rows2.Close()
	for rows2.Next() {
		var tbl, cname, ctype, col string
		var ftbl, fcol sql.NullString
		rows2.Scan(&tbl, &cname, &ctype, &col, &ftbl, &fcol)
		if ctype == "FOREIGN KEY" {
			fmt.Printf("  %s.%s → %s.%s (%s)\n", tbl, col, ftbl.String, fcol.String, cname)
		} else {
			fmt.Printf("  %s.%s [%s] (%s)\n", tbl, col, ctype, cname)
		}
	}
}

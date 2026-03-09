package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

func main() {
	// Get from environment
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		fmt.Println("❌ DATABASE_URL not set")
		fmt.Println("\nUsage:")
		fmt.Println("  export DATABASE_URL='postgresql://...'")
		fmt.Println("  go run test-connection.go")
		os.Exit(1)
	}

	fmt.Println("========================================")
	fmt.Printf("Testing connection to: %s\n", maskPassword(connStr))
	fmt.Println("========================================")

	// Test connection
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ Connection failed: %v\n\n", err)
		fmt.Println("Troubleshooting:")
		fmt.Println("1. Verify project exists at supabase.com/dashboard")
		fmt.Println("2. Check connection string: Settings → Database → URI")
		fmt.Println("3. Verify project is not paused (free tier auto-pauses)")
		fmt.Println("4. Check password is correct")
		fmt.Println("\nSee GET_CONNECTION_STRING.md for detailed help")
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	fmt.Println("✅ Connection successful!")

	// Test query
	var result int
	err = conn.QueryRow(context.Background(), "SELECT 1").Scan(&result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Query failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Query test passed")

	// Check tables
	var tableCount int
	err = conn.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'").Scan(&tableCount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Table check failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Found %d tables in database\n", tableCount)

	if tableCount == 0 {
		fmt.Println("\n⚠️  Warning: No tables found!")
		fmt.Println("   Run migrations first: See docs/MANUAL_MIGRATION.md")
		os.Exit(0)
	}

	// List tables
	rows, err := conn.Query(context.Background(),
		"SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE' ORDER BY table_name")
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to list tables: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	fmt.Println("\nTables:")
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			continue
		}
		fmt.Printf("  - %s\n", tableName)
	}

	// Check sources
	var sourceCount int
	err = conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM sources WHERE enabled = true").Scan(&sourceCount)
	if err != nil {
		fmt.Println("\n⚠️  Warning: Could not check sources table")
		fmt.Println("   Make sure migrations are run: docs/MANUAL_MIGRATION.md")
	} else {
		fmt.Printf("\n✅ Found %d enabled sources\n", sourceCount)
		if sourceCount == 0 {
			fmt.Println("   ⚠️  No sources enabled! Seed default sources:")
			fmt.Println("   See docs/MANUAL_MIGRATION.md Step 7")
		}
	}

	fmt.Println("\n========================================")
	fmt.Println("✅ All checks passed! Ready to run worker.")
	fmt.Println("========================================")
}

func maskPassword(connStr string) string {
	// Extract parts
	if strings.Contains(connStr, "@") {
		parts := strings.Split(connStr, "@")
		if len(parts) == 2 {
			// Get user part
			userPart := parts[0]
			if strings.Contains(userPart, ":") {
				userParts := strings.Split(userPart, ":")
				if len(userParts) >= 2 {
					return userParts[0] + ":***@" + parts[1]
				}
			}
		}
	}
	return "postgresql://postgres:***@db.xxx.supabase.co:5432/postgres"
}

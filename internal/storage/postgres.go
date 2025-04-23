package storage

import (
	"database/sql"
	"fmt"
	// "os"
	// "strings"

	// "github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitPostgres() error {
	connStr := "postgres://user:password@postgres:5432/finTechDB?sslmode=disable"
	// godotenv.Load()
	// connStr := strings.TrimSpace(os.Getenv("DB_URL"))
	
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open Postgres connection: %w", err)
	}

	return DB.Ping()
}
package persistence

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // postgres driver
)

func ConnectToPostgres() (*sql.DB, error) {
	host := os.Getenv("DB_HOSTNAME")
	if host == "" {
		host = "localhost"
	}
	connStr := "host=" + host + " user=admin password=123 dbname=rinha sslmode=disable"
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("open %w", err)
	}
	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("ping %w", err)
	}
	return database, nil
}

package persistence

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq" // postgres driver
)

func ConnectToPostgres(log func(string, ...any)) (*sql.DB, error) {
	log(`Connecting to Postgres...`)
	defer log(`Connecting to Postgres...(done)`)
	host := os.Getenv("DB_HOSTNAME")
	if host == "" {
		host = "localhost"
	}
	connStr := "host=" + host + " user=admin password=123 dbname=rinha sslmode=disable"
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("open %w", err)
	}
	var pingErr error
	for i := 0; ; i++ {
		if pingErr = database.Ping(); pingErr == nil {
			break
		}
		time.Sleep(time.Second * time.Duration(i))
	}
	if pingErr != nil {
		return nil, fmt.Errorf("ping %w", err)
	}
	return database, nil
}

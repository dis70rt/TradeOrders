package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	log "github.com/dis70rt/TradeOrders/internals/logger"
)

func ConnectPostgres() (*sql.DB) {
	_ = godotenv.Load()
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.WithError(err).Fatal("Failed to open database connection")
		return nil
	}
	
	if err := db.Ping(); err != nil {
		log.WithError(err).Fatal("Failed to connect to the database")
		return nil
	}
	
	return db
}
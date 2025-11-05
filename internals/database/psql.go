package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	log "github.com/dis70rt/TradeOrders/internals/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)
    db.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		log.WithError(err).Fatal("Failed to connect to the database")
		return nil
	}
	
	return db
}
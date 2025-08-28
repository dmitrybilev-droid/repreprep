package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewTestDB() *sql.DB {
	err := godotenv.Load("pkg/db/.env.test")
	if err != nil {
		log.Fatalf("Failed to load test env: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to test DB: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Test DB is not reachable: %v", err)
	}

	return db
}

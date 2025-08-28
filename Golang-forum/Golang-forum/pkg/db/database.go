package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"golang-forum/pkg/config"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB ping failed: %v", err)
	}

	return db, nil
}

package db

import (
	"beli-mang/internal/pkg/configuration"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func New(cfg *configuration.Configuration) (*sql.DB, error) {
	// connect to PostgreSQL
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?%v",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBParams)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// check connection
	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db, nil
}

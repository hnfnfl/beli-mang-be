package db

import (
	"beli-mang/internal/pkg/configuration"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

func GetConn(cfg *configuration.Configuration, ctx context.Context) *pgxpool.Pool {
	// connect to PostgreSQL
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?%v",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBParams)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	config.MaxConnLifetime = 60 * time.Minute
	config.MaxConnIdleTime = time.Minute * 10
	config.ConnConfig.ConnectTimeout = time.Second * 5

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	log.Println("success connect to the database")

	return pool
}

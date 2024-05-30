package db

import (
	"beli-mang/internal/pkg/configuration"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/jackc/pgx/v5/tracelog"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func GetConn(ctx context.Context, cfg *configuration.Configuration, log *logrus.Logger) *pgxpool.Pool {
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

	// NOTE: Uncomment this if you want to use tracelog, but not recommended for production
	// tracer := &tracelog.TraceLog{
	// 	LogLevel: tracelog.LogLevel(log.Level),
	// 	Logger: tracelog.LoggerFunc(
	// 		func(ctx context.Context, level tracelog.LogLevel, _ string, data map[string]interface{}) {
	// 			log.WithFields(data).Log(log.Level)
	// 		},
	// 	),
	// }
	// config.ConnConfig.Tracer = tracer

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	log.Infof("Connected to DB: %v", cfg.DBName)

	return pool
}

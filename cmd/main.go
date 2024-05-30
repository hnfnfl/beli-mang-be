package main

import (
	"beli-mang/internal/pkg/configuration"
	"beli-mang/internal/pkg/logger"
	"beli-mang/internal/pkg/server"
)

func main() {
	config, err := configuration.NewConfiguration()
	if err != nil {
		panic(err)
	}

	log, err := logger.NewLogger(config.LogLevel)
	if err != nil {
		panic(err)
	}

	if err := server.Run(config, log); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

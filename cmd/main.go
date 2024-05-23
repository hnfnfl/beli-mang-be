package main

import (
	"beli-mang/internal/pkg/configuration"
	"beli-mang/internal/pkg/handler"
	"beli-mang/internal/pkg/logger"
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

	if err := handler.Run(config, log); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

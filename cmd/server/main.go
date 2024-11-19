package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arefev/mtrcstore/internal/server"
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/arefev/mtrcstore/internal/server/logger"
	"github.com/arefev/mtrcstore/internal/server/repository"
	"github.com/arefev/mtrcstore/internal/server/worker"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := NewConfig(os.Args[1:])

	if err != nil {
		return fmt.Errorf("main config init failed: %w", err)
	}

	cLog, err := logger.Build(config.LogLevel)
	if err != nil {
		return fmt.Errorf("logger init failed: %w", err)
	}

	storage := repository.NewMemory()
	metricHandlers := handler.NewMetricHandlers(&storage, cLog)

	go worker.
		Init(config.StoreInterval, config.FileStoragePath, config.Restore, &storage, cLog).
		Run()

	r := server.InitRouter(metricHandlers, cLog)

	log.Printf("Server up on address %s\n", config.Address)
	log.Printf("Log level %s\n", config.LogLevel)
	return fmt.Errorf("main run() failed: %w", http.ListenAndServe(config.Address, r))
}

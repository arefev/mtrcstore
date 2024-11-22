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
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
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

	cLog.Info("filePath", zap.Int("", len(config.FileStoragePath)))

	storage, storageType, err := initStorage(&config, cLog)
	if err != nil {
		return fmt.Errorf("main run failed: %w", err)
	}

	metricHandlers := handler.NewMetricHandlers(storage, cLog)
	r := server.InitRouter(metricHandlers, cLog)

	cLog.Info(
		"Server running",
		zap.String("address", config.Address),
		zap.String("log level", config.LogLevel),
		zap.String("storage type", storageType),
	)

	return fmt.Errorf("main run() failed: %w", http.ListenAndServe(config.Address, r))
}

func initStorage(config *Config, cLog *zap.Logger) (storage repository.Storage, storageType string, err error) {
	switch {
	case len(config.DatabaseDSN) > 0:
		storage, err = repository.NewDatabaseRep(config.DatabaseDSN, cLog)
		if err != nil {
			err = fmt.Errorf("repository init failed: %w", err)
		}

		storageType = "DB"
	case len(config.FileStoragePath) > 0:
		storage = repository.
			NewFile(config.StoreInterval, config.FileStoragePath, config.Restore, cLog).
			WorkerRun()

		storageType = "File"
	default:
		storage = repository.NewMemory()
		storageType = "Memory"
	}

	return
}

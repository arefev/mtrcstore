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

	// storage := repository.
	// 	NewFile(config.StoreInterval, config.FileStoragePath, config.Restore, cLog).
	// 	WorkerRun()

	storage, err := repository.NewDatabaseRep(config.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("repository init failed: %w", err)
	}

	metricHandlers := handler.NewMetricHandlers(storage, cLog)
	r := server.InitRouter(metricHandlers, cLog)

	cLog.Info(
		"Server running",
		zap.String("address", config.Address),
		zap.String("log level", config.LogLevel),
	)

	return fmt.Errorf("main run() failed: %w", http.ListenAndServe(config.Address, r))
}

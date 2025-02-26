package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arefev/mtrcstore/internal/agent"
	"github.com/arefev/mtrcstore/internal/agent/repository"
	"github.com/arefev/mtrcstore/internal/agent/service"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	config, err := NewConfig(os.Args[1:])

	if err != nil {
		return fmt.Errorf("main config init failed: %w", err)
	}

	requestClient := service.Client{}

	storage := repository.NewMemory()
	report := service.NewReport(&storage, config.Address, config.SecretKey, requestClient)

	worker := agent.Worker{
		WorkerPool:     service.NewWorkerPool(report, config.RateLimit),
		PollInterval:   config.PollInterval,
		ReportInterval: config.ReportInterval,
	}

	log.Printf(
		"Run worker with params:\nserverHost = %s\npollInterval = %d\nreportInterval = %d\nrateLimit = %d\n",
		config.Address,
		config.PollInterval,
		config.ReportInterval,
		config.RateLimit,
	)

	return fmt.Errorf("main run() failed: %w", worker.Run())
}

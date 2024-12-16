package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arefev/mtrcstore/internal/agent"
	"github.com/arefev/mtrcstore/internal/agent/repository"
	"github.com/arefev/mtrcstore/internal/agent/service"
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

	storage := repository.NewMemory()
	report, err := service.NewReport(&storage, config.Address, config.SecretKey)
	if err != nil {
		return fmt.Errorf("main run() failed: %w", err)
	}

	worker := agent.Worker{
		WorkerPool:     service.NewWorkerPool(&report, config.RateLimit),
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

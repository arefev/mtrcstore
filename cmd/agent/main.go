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
	report, err := service.NewReport(&storage, config.Address)
	if err != nil {
		return fmt.Errorf("main run() failed: %w", err)
	}

	worker := agent.Worker{
		Report:         &report,
		PollInterval:   config.PollInterval,
		ReportInterval: config.ReportInterval,
	}

	log.Printf(
		"Run worker with params:\nserverHost = %s\npollInterval = %d\nreportInterval = %d",
		config.Address,
		config.PollInterval,
		config.ReportInterval,
	)

	return fmt.Errorf("main run() failed: %w", worker.Run())
}

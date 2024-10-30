package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arefev/mtrcstore/internal/agent"
	"github.com/arefev/mtrcstore/internal/agent/repository"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config := NewConfig()

	if err := config.InitFlags(os.Args[1:]); err != nil {
		return err
	}

	if err := config.InitEnvs(); err != nil {
		return err
	}

	storage := repository.NewMemory()
	worker := agent.Worker{
		PollInterval:   config.PollInterval,
		ReportInterval: config.ReportInterval,
		Storage:        &storage,
		ServerHost:     config.Address,
	}

	log.Printf(
		"Run worker with params:\nserverHost = %s\npollInterval = %d\nreportInterval = %d",
		config.Address,
		config.PollInterval,
		config.ReportInterval,
	)

	return fmt.Errorf("main run() failed: %w", worker.Run())
}

package main

import (
	"context"
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
	ctx := context.Background()
	requestClient := service.Client{}
	if err := run(ctx, os.Args[1:], &requestClient); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, args []string, sender service.Sender) error {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	config, err := NewConfig(args)
	if err != nil {
		return fmt.Errorf("main config init failed: %w", err)
	}

	storage := repository.NewMemory()
	report := service.NewReport(&storage, config.Address, config.SecretKey, sender)

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

	return fmt.Errorf("main run() failed: %w", worker.Run(ctx))
}

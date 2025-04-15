package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	config, err := NewConfig(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	var client service.Sender
	switch {
	case config.GRPCAddress != "":
		client = service.NewGRPCClient(config.GRPCAddress)
	default:
		client = service.NewClient(
			config.SecretKey,
			config.CryptoKey,
			"http://"+config.Address+"/updates/",
		)
	}

	if err := run(ctx, &config, client); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, config *Config, sender service.Sender) error {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	storage := repository.NewMemory()
	report := service.NewReport(&storage, sender)

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

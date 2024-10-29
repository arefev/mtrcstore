package main

import (
	"log"

	"github.com/arefev/mtrcstore/internal/agent"
	"github.com/arefev/mtrcstore/internal/agent/repository"
)

func main() {
	config := NewConfig()

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
	worker.Run()
}

package main

import (
	"github.com/arefev/mtrcstore/internal/agent"
	"github.com/arefev/mtrcstore/internal/agent/repository"
)

func main() {
	config := ParseFlags()

	storage := repository.NewMemory()
	worker := agent.Worker{
		PollInterval:   config.PollInterval,
		ReportInterval: config.ReportInterval,
		Storage:        &storage,
		ServerHost:     config.Address,
	}
	worker.Run()
}

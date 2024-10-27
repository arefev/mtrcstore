package main

import (
	"github.com/arefev/mtrcstore/internal/agent"
	"github.com/arefev/mtrcstore/internal/agent/repository"
)

func main() {
	parseFlags()

	storage := repository.NewMemory()
	worker := agent.Worker{
		PollInterval:   flagPollInterval,
		ReportInterval: flagReportInterval,
		Storage:        &storage,
		ServerHost:     flagServerAddr,
	}
	worker.Run()
}

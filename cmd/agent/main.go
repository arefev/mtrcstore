package main

import (
	"github.com/arefev/mtrcstore/internal/agent"
	"github.com/arefev/mtrcstore/internal/agent/repository"
)


func main() {
	const pollInterval = 2
	const reportInterval = 10
	const serverHost = "http://localhost:8080"

	storage := repository.NewMemory()
	worker := agent.Worker{
		PollInterval: pollInterval,
		ReportInterval: reportInterval,
		Storage: &storage,
		ServerHost: serverHost,
	}
	worker.Run()
}
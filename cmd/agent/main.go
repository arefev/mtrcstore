package main

import (
	"github.com/arefev/mtrcstore/internal/agent"
	"github.com/arefev/mtrcstore/internal/agent/repository"
)


func main() {
	const pollInterval = 2
	const reportInterval = 10

	storage := &repository.Memory{}
	worker := agent.Worker{
		PollInterval: pollInterval,
		ReportInterval: reportInterval,
		Storage: storage,
	}
	worker.Run()
}
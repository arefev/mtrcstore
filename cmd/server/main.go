package main

import (
	"log"
	"net/http"
	"os"

	"github.com/arefev/mtrcstore/internal/server"
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/arefev/mtrcstore/internal/server/repository"
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
	metricHandlers := handler.MetricHandlers{
		Storage: &storage,
	}

	r := server.InitRouter(&metricHandlers)

	log.Printf("Server up on address %s\n", config.Address)
	return http.ListenAndServe(config.Address, r)
}

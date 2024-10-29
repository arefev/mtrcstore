package main

import (
	"log"
	"net/http"

	"github.com/arefev/mtrcstore/internal/server"
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/arefev/mtrcstore/internal/server/repository"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := run(config); err != nil {
		log.Fatal(err)
	}
}

func run(config Config) error {
	storage := repository.NewMemory()
	handler := handler.MetricHandlers{
		Storage: &storage,
	}

	r := server.InitRouter(&handler)

	log.Printf("Server up on address %s\n", config.Address)
	return http.ListenAndServe(config.Address, r)
}

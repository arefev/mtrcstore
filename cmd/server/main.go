package main

import (
	"log"
	"net/http"

	"github.com/arefev/mtrcstore/internal/server"
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/arefev/mtrcstore/internal/server/repository"
)

func main() {
	config := NewConfig()
	if err := run(config); err != nil {
		panic(err)
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

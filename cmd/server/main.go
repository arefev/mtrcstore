package main

import (
	"fmt"
	"net/http"

	"github.com/arefev/mtrcstore/internal/server"
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/arefev/mtrcstore/internal/server/repository"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	const addr = "localhost:8080"

	storage := repository.NewMemory()
	handler := handler.MetricHandlers{
		Storage: &storage,
	}

	r := server.InitRouter(&handler)

	fmt.Printf("Server up by address %s\n", addr)
	return http.ListenAndServe(addr, r)
}

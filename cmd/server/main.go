package main

import (
	"fmt"
	"net/http"

	"github.com/arefev/mtrcstore/internal/server"
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/arefev/mtrcstore/internal/server/repository"
)

func main() {
	parseFlags()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	storage := repository.NewMemory()
	handler := handler.MetricHandlers{
		Storage: &storage,
	}

	r := server.InitRouter(&handler)

	fmt.Printf("Server up by address %s\n", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, r)
}

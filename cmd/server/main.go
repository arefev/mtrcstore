package main

import (
	"fmt"
	"net/http"

	"github.com/arefev/mtrcstore/internal/server/http/handler"
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
	handler := handler.UpdateHandler{
		Storage: &storage,
	}

	mux := http.NewServeMux()
	handler.Handle(mux)

	fmt.Printf("Server up by address %s\n", addr)
	return http.ListenAndServe(addr, mux)
}

package main

import (
	"fmt"
	"net/http"

	"github.com/arefev/mtrcstore/internal/server/http/handler"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	const addr = "localhost:8080"

	mux := http.NewServeMux()
	handler.UpdateHandler(mux)

	fmt.Printf("Server up by address %s\n", addr)
	return http.ListenAndServe(addr, mux)
}

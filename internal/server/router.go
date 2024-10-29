package server

import (
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouter(handler *handler.MetricHandlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", handler.Get)
	r.Get("/value/{type}/{name}", handler.Find)
	r.Post("/update/{type}/{name}/{value}", handler.Update)

	return r
}

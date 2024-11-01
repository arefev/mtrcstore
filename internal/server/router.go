package server

import (
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouter(h *handler.MetricHandlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", h.Get)
	r.Get("/value/{type}/{name}", h.Find)
	r.Post("/update/{type}/{name}/{value}", h.Update)

	return r
}

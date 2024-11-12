package server

import (
	"github.com/arefev/mtrcstore/internal/server/handler"
	// "github.com/arefev/mtrcstore/internal/server/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRouter(h *handler.MetricHandlers) *chi.Mux {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)

	r.Get("/", h.Get)
	r.Get("/value/{type}/{name}", h.Find)
	r.Post("/update/{type}/{name}/{value}", h.Update)
	r.Post("/value", h.FindJson)
	r.Post("/update", h.UpdateJson)

	return r
}

package server

import (
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/arefev/mtrcstore/internal/server/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func InitRouter(handler *handler.MetricHandlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)

	r.Get("/", handler.Get)

	r.Route("/value/{type}", func(r chi.Router) {
		r.Use(middleware.CheckType)
		r.Get("/{name}", handler.Find)
	})

	r.Route("/update/{type}", func(r chi.Router) {
		r.Use(middleware.CheckType)
		r.Post("/{name}/{value}", handler.Update)
	})

	return r
}

package server

import (
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/arefev/mtrcstore/internal/server/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func InitRouter(h *handler.MetricHandlers, log *zap.Logger) *chi.Mux {
	m := middleware.NewMiddleware(log)
	r := chi.NewRouter()	
	r.Use(m.Logger, m.Compress)

	r.Get("/", h.Get)

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", h.Find)
		r.Post("/", h.FindJSON)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", h.Update)
		r.Post("/", h.UpdateJSON)
	})

	return r
}

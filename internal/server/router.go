package server

import (
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/arefev/mtrcstore/internal/server/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

func InitRouter(h *handler.MetricHandlers, log *zap.Logger, secretKey string, cryptoKey string) *chi.Mux {
	m := middleware.NewMiddleware(log, secretKey, cryptoKey)
	r := chi.NewRouter()
	r.Use(m.Logger)
	r.Use(m.Compress)
	r.Use(m.Decrypt)
	r.Use(m.CheckSign)
	r.Mount("/debug", chi_middleware.Profiler())

	r.Get("/", h.Get)
	r.Get("/ping", h.Ping)

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", h.Find)
		r.Post("/", h.FindJSON)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", h.Update)
		r.Post("/", h.UpdateJSON)
	})

	r.Post("/updates/", h.Updates)

	return r
}

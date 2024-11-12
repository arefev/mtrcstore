package middleware

import (
	"net/http"

	"github.com/arefev/mtrcstore/internal/server/logger"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Info(r.Method)
		next.ServeHTTP(w, r)
	})
}
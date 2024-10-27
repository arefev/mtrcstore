package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CheckType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		mType := chi.URLParam(r, "type")
		if mType != "counter" && mType != "gauge" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

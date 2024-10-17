package middleware

import "net/http"

func PostOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.Method != http.MethodPost) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

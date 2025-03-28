package middleware

import (
	"net"
	"net/http"
)

func (m *Middleware) IsPrivateIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.cidr == "" {
			next.ServeHTTP(w, r)
			return
		}

		ip := net.ParseIP(r.Header.Get("X-Real-IP"))
		if ip == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		_, block, err := net.ParseCIDR(m.cidr)
		if err != nil || !block.Contains(ip) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

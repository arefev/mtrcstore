package middleware

import (
	"net/http"
	"strings"

	"github.com/arefev/mtrcstore/internal/server/service"
	"go.uber.org/zap"
)

func (m *Middleware) Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Accept")
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		if supportsGzip && checkContentType(contentType) {
			ow := service.NewCompressWriter(w)
			w = ow
			defer func() {
				if err := ow.Close(); err != nil {
					m.log.Debug("writer body close error", zap.Error(err))
				}
			}()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			cr, err := service.NewCompressReader(r.Body)
			if err != nil {
				m.log.Debug("gzip error", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer func() {
				if err := cr.Close(); err != nil {
					m.log.Debug("reader body close error", zap.Error(err))
				}
			}()
		}

		next.ServeHTTP(w, r)
	})
}

func checkContentType(cType string) bool {
	return strings.Contains(cType, "application/json") || strings.Contains(cType, "text/html")
}

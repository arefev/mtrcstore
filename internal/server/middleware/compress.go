package middleware

import (
	"net/http"
	"strings"

	"github.com/arefev/mtrcstore/internal/server/service"
	"go.uber.org/zap"
)

func (m *Middleware) Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		contentType := r.Header.Get("Accept")
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		if supportsGzip && checkContentType(contentType) {
			cw := service.NewCompressWriter(w)
			ow = cw
			defer func() {
				if err := cw.Close(); err != nil {
					m.log.Debug("writer body close error", zap.Error(err))
				}
			}()

			// Почему-то не вызывается WriteHeader из ow
			ow.Header().Set("Content-Encoding", "gzip")
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

		next.ServeHTTP(ow, r)
	})
}

func checkContentType(cType string) bool {
	return strings.Contains(cType, "application/json") || strings.Contains(cType, "text/html")
}

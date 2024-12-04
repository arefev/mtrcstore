package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
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
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				m.log.Debug("gzip reader error", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			b, err := io.ReadAll(gz)
			if err != nil {
				m.log.Debug("gzip reader error", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			defer func() {
				if err := gz.Close(); err != nil {
					m.log.Debug("gzip reader error", zap.Error(err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}()

			r.Body = io.NopCloser(bytes.NewBuffer(b))
		}

		next.ServeHTTP(w, r)
	})
}

func checkContentType(cType string) bool {
	return strings.Contains(cType, "application/json") || strings.Contains(cType, "text/html")
}

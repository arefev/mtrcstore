package middleware

import (
	"net/http"
	"strings"

	"github.com/arefev/mtrcstore/internal/server/logger"
	"github.com/arefev/mtrcstore/internal/server/service"
	"go.uber.org/zap"
)

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := service.NewCompressReader(r.Body)
			if err != nil {
				logger.Log.Debug("gzip error", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
			}
			r.Body = cr
			defer func() {
				if err := cr.Close(); err != nil {
					logger.Log.Debug("body close error", zap.Error(err))
				}
			}()
		}

		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)

	if err != nil {
		return size, fmt.Errorf("loggingResponseWriter Write: %w", err)
	}

	r.responseData.size += size
	return size, nil
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: http.StatusOK,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		m.log.Info(
			"Request handler",
			zap.String("URI", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", responseData.status),
			zap.Duration("duration", duration),
			zap.Int("size", responseData.size),
		)
	})
}

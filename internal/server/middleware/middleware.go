package middleware

import "go.uber.org/zap"

type Middleware struct {
	log *zap.Logger
}

func NewMiddleware(log *zap.Logger) Middleware {
	return Middleware{
		log: log,
	}
}

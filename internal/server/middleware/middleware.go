package middleware

import "go.uber.org/zap"

type Middleware struct {
	secretKey string
	log       *zap.Logger
}

func NewMiddleware(log *zap.Logger, secretKey string) Middleware {
	return Middleware{
		log:       log,
		secretKey: secretKey,
	}
}

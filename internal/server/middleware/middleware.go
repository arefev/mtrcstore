package middleware

import "go.uber.org/zap"

type Middleware struct {
	log       *zap.Logger
	secretKey string
	cryptoKey string
}

func NewMiddleware(log *zap.Logger, secretKey string, cryptoKey string) Middleware {
	return Middleware{
		log:       log,
		secretKey: secretKey,
		cryptoKey: cryptoKey,
	}
}

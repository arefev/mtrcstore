package middleware

import "go.uber.org/zap"

type Middleware struct {
	log       *zap.Logger
	secretKey string
	cryptoKey string
	cidr      string
}

func NewMiddleware(log *zap.Logger, cidr string, secretKey string, cryptoKey string) Middleware {
	return Middleware{
		log:       log,
		cidr:      cidr,
		secretKey: secretKey,
		cryptoKey: cryptoKey,
	}
}

package middleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

func (m *Middleware) Decrypt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.cryptoKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			m.log.Error("middleware Decrypt: read body failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err = decrypt(body, m.cryptoKey)
		if err != nil {
			m.log.Error("middleware Decrypt: decrypt failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	})
}

func decrypt(data []byte, cryptoKey string) ([]byte, error) {
	privateKeyPEM, err := os.ReadFile(cryptoKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt - ReadFile with private key failed: %w", err)
	}

	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("decrypt - ParsePKCS1PrivateKey failed: %w", err)
	}

	msgLen := len(data)
	step := privateKey.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := min(start+step, msgLen)

		decryptedBlockBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data[start:finish])
		if err != nil {
			return nil, fmt.Errorf("decrypt - DecryptPKCS1v15 failed: %w", err)
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (m *Middleware) CheckSign(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey := []byte("312msdlfmaskn1223lmn123ns")
		if len(secretKey) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		hash := r.Header.Get("HashSHA256")
		if hash == "" {
			m.log.Error("check sign failed: hash is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			m.log.Error("check sign failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		copy := io.NopCloser(bytes.NewBuffer(body))
		sign, err := sign(secretKey, body)
		if err != nil {
			m.log.Error("check sign failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = copy

		hashDecoded, err := hex.DecodeString(hash)
		if err != nil {
			m.log.Error("check sign failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !hmac.Equal(sign, hashDecoded) {
			m.log.Error("check sign failed: hashs not equal")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func sign(secretKey []byte, data []byte) ([]byte, error) {
	h := hmac.New(sha256.New, secretKey)

	if _, err := h.Write(data); err != nil {
		return []byte{}, fmt.Errorf("sign failed: %w", err)
	}

	dst := h.Sum(nil)
	return dst, nil
}

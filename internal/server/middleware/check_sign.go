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

type signWriter struct {
	http.ResponseWriter
	secretKey []byte
}

func NewSignWriter(w http.ResponseWriter, secretKey []byte) *signWriter {
	return &signWriter{
		ResponseWriter: w,
		secretKey:      secretKey,
	}
}

func (s *signWriter) Write(p []byte) (int, error) {
	hash, err := sign(s.secretKey, p)
	if err != nil {
		return 0, fmt.Errorf("write failed: %w", err)
	}

	s.Header().Add("HashSHA256", hex.EncodeToString(hash))
	n, err := s.ResponseWriter.Write(p)
	if err != nil {
		return 0, fmt.Errorf("write failed: %w", err)
	}

	return n, nil
}

func (m *Middleware) CheckSign(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey := []byte(m.secretKey)
		hash := r.Header.Get("HashSHA256")

		if len(secretKey) == 0 || hash == "" {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			m.log.Error("check sign failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bodyCopy := io.NopCloser(bytes.NewBuffer(body))
		sign, err := sign(secretKey, body)
		if err != nil {
			m.log.Error("check sign failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = bodyCopy

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

		newResp := NewSignWriter(w, secretKey)
		w = newResp

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

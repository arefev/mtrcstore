package service

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/arefev/mtrcstore/internal/agent/model"
	"github.com/go-resty/resty/v2"
)

var ErrRequestFail = errors.New("doRequest failed")

type client struct {
	secretKey string
	cryptoKey string
	url       string
}

func NewClient(secretKey, cryptoKey, url string) *client {
	return &client{
		secretKey: secretKey,
		cryptoKey: cryptoKey,
		url:       url,
	}
}

func (c *client) doRequest(ctx context.Context, headers map[string]string, body any) error {
	request := resty.New().R().SetContext(ctx)
	for k, v := range headers {
		request.SetHeader(k, v)
	}

	if _, err := request.SetBody(body).Post(c.url); err != nil {
		return fmt.Errorf("%w: %w", ErrRequestFail, err)
	}

	return nil
}

func (c *client) Request(ctx context.Context, data []model.Metric) error {
	headers := map[string]string{
		"Content-Type":     "application/json",
		"Content-Encoding": "gzip",
	}
	
	jsonBody, mErr := json.Marshal(data)
	if mErr != nil {
		return c.requestError(mErr)
	}

	if c.secretKey != "" {
		hash, err := c.sign(jsonBody)
		if err != nil {
			return c.requestError(err)
		}

		headers["HashSHA256"] = hex.EncodeToString(hash)
	}

	ip, err := c.getIP()
	if err != nil {
		return c.requestError(err)
	}
	headers["X-Real-IP"] = ip

	body, err := c.compress(jsonBody)
	if err != nil {
		return c.requestError(err)
	}

	jsonBody, err = io.ReadAll(body)
	if err != nil {
		return c.requestError(err)
	}

	if c.cryptoKey != "" {
		ecrypted, err := c.encrypt(jsonBody, c.cryptoKey)
		if err != nil {
			return c.requestError(err)
		}

		jsonBody = ecrypted
	}

	if err := c.doRequest(ctx, headers, jsonBody); err != nil {
		return c.requestError(err)
	}

	return nil
}

func (c *client) getIP() (string, error) {
	var ip net.IP
	a, _ := net.Interfaces()
	for _, i := range a {
		addrs, err := i.Addrs()
		if err != nil {
			return "", fmt.Errorf("getIP addrs failed: %w", err)
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
		}
	}

	if ip == nil {
		return "", errors.New("getIP failed: no IP address found")
	}

	return ip.String(), nil
}

func (c *client) requestError(err error) error {
	return fmt.Errorf("request failed: %w", err)
}

func (c *client) encrypt(data []byte, cryptoKey string) ([]byte, error) {
	const decreaseKeySize int = 11

	publicKeyPEM, err := os.ReadFile(cryptoKey)
	if err != nil {
		return nil, fmt.Errorf("encrypt - ReadFile failed: %w", err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	parsed, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("encrypt - Decode failed: %w", err)
	}

	publicKey, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("encrypt - invalid public key type")
	}

	msgLen := len(data)
	step := publicKey.Size() - decreaseKeySize
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := min(start+step, msgLen)
		encryptedBlockBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data[start:finish])
		if err != nil {
			return nil, fmt.Errorf("encrypt - EncryptPKCS1v15 failed: %w", err)
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func (c *client) sign(data []byte) ([]byte, error) {
	key := []byte(c.secretKey)
	h := hmac.New(sha256.New, key)

	if _, err := h.Write(data); err != nil {
		return []byte{}, fmt.Errorf("sign failed: %w", err)
	}

	dst := h.Sum(nil)
	return dst, nil
}

func (c *client) compress(p []byte) (*bytes.Buffer, error) {
	var err error
	body := bytes.NewBuffer(nil)
	w := gzip.NewWriter(body)
	if _, err = w.Write(p); err != nil {
		return body, fmt.Errorf("gzip failed: %w", err)
	}

	if err := w.Close(); err != nil {
		return body, fmt.Errorf("gzip failed: %w", err)
	}

	return body, nil
}

func (c *client) IsConnRefused(err error) bool {
	var netErr *net.OpError
	return errors.As(err, &netErr) && netErr.Op == "dial" && netErr.Err.Error() == "connect: connection refused"
}

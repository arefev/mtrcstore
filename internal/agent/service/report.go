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
	"log"
	"net"
	"os"
	"runtime"

	"github.com/arefev/mtrcstore/internal/agent/model"
	"github.com/arefev/mtrcstore/internal/proto"
	// "github.com/arefev/mtrcstore/internal/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gauge float64
type Counter int64

type Storage interface {
	Save(memStats *runtime.MemStats) error
	SaveCPU() error
	IncrementCounter()
	ClearCounter()
	GetGauges() map[string]Gauge
	GetCounters() map[string]Counter
}

type Sender interface {
	DoRequest(ctx context.Context, url string, headers map[string]string, body any) error
}

type Report struct {
	Storage       Storage
	sender        Sender
	updateURL     string
	massUpdateURL string
	gaugeName     string
	counterName   string
	secretKey     string
	host          string
	cryptoKey     string
}

func NewReport(s Storage, host string, secretKey string, cryptoKey string, sender Sender) *Report {
	const (
		protocol          = "http://"
		updateURLPath     = "/update"
		massUpdateURLPath = "/updates/"
		counterName       = "counter"
		gaugeName         = "gauge"
	)

	return &Report{
		Storage:       s,
		updateURL:     updateURLPath,
		massUpdateURL: massUpdateURLPath,
		gaugeName:     gaugeName,
		counterName:   counterName,
		secretKey:     secretKey,
		cryptoKey:     cryptoKey,
		sender:        sender,
		host:          protocol + host,
	}
}

func (r *Report) Send(ctx context.Context, metrics []model.Metric) {
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := proto.NewMetricsClient(conn)

	var pMetrics []*proto.Metric
	for _, m := range metrics {
		pm := &proto.Metric{
			ID:   m.ID,
			Type: m.MType,
		}

		if m.Value != nil {
			pm.Value = *m.Value
		}

		if m.Delta != nil {
			pm.Delta = *m.Delta
		}

		pMetrics = append(pMetrics, pm)
	}

	response, err := c.UpdateMetric(context.Background(), &proto.UpdateMetricRequest{
		Metrics: pMetrics,
	})

	if err != nil {
		log.Printf("send(): failed to send the metrics: %s", err.Error())
		return
	}

	log.Printf("response: %+v", response)

	// const rCount = 3
	// action := func() error {
	// 	return r.request(ctx, metrics, r.massUpdateURL)
	// }
	// if err := retry.New(action, r.isConnRefused, rCount).Run(); err != nil {
	// 	log.Printf("sendCounters(): failed to send the counter metric %s: %s", r.counterName, err.Error())
	// }
}

func (r *Report) GetMetrics() []model.Metric {
	metrics := make([]model.Metric, 0)
	metrics = append(metrics, r.getGauges()...)
	metrics = append(metrics, r.getCounters()...)
	return metrics
}

func (r *Report) getGauges() []model.Metric {
	metrics := make([]model.Metric, 0)
	for name, val := range r.Storage.GetGauges() {
		mVal := float64(val)
		metrics = append(metrics, model.Metric{
			ID:    name,
			MType: r.gaugeName,
			Value: &mVal,
		})
	}

	return metrics
}

func (r *Report) getCounters() []model.Metric {
	metrics := make([]model.Metric, 0)
	for name, val := range r.Storage.GetCounters() {
		delta := int64(val)
		metrics = append(metrics, model.Metric{
			ID:    name,
			MType: r.counterName,
			Delta: &delta,
		})
	}

	return metrics
}

func (r *Report) Save(memStats *runtime.MemStats) error {
	if err := r.Storage.Save(memStats); err != nil {
		return fmt.Errorf("report save(): metrics save failed: %w", err)
	}

	return nil
}

func (r *Report) SaveCPU() error {
	if err := r.Storage.SaveCPU(); err != nil {
		return fmt.Errorf("report saveCPU() failed: %w", err)
	}

	return nil
}

func (r *Report) IncrementCounter() {
	r.Storage.IncrementCounter()
}

func (r *Report) ClearCounter() {
	r.Storage.ClearCounter()
}

func (r *Report) request(ctx context.Context, data any, url string) error {
	headers := map[string]string{
		"Content-Type":     "application/json",
		"Content-Encoding": "gzip",
	}

	jsonBody, mErr := json.Marshal(data)
	if mErr != nil {
		return r.requestError(mErr)
	}

	if r.secretKey != "" {
		hash, err := r.sign(jsonBody)
		if err != nil {
			return r.requestError(err)
		}

		headers["HashSHA256"] = hex.EncodeToString(hash)
	}

	ip, err := r.getIP()
	if err != nil {
		return r.requestError(err)
	}
	headers["X-Real-IP"] = ip

	body, err := r.compress(jsonBody)
	if err != nil {
		return r.requestError(err)
	}

	jsonBody, err = io.ReadAll(body)
	if err != nil {
		return r.requestError(err)
	}

	if r.cryptoKey != "" {
		ecrypted, err := r.encrypt(jsonBody, r.cryptoKey)
		if err != nil {
			return r.requestError(err)
		}

		jsonBody = ecrypted
	}

	if err := r.sender.DoRequest(ctx, r.host+url, headers, jsonBody); err != nil {
		return r.requestError(err)
	}

	return nil
}

func (r *Report) getIP() (string, error) {
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

func (r *Report) requestError(err error) error {
	return fmt.Errorf("request failed: %w", err)
}

func (r *Report) encrypt(data []byte, cryptoKey string) ([]byte, error) {
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

func (r *Report) sign(data []byte) ([]byte, error) {
	key := []byte(r.secretKey)
	h := hmac.New(sha256.New, key)

	if _, err := h.Write(data); err != nil {
		return []byte{}, fmt.Errorf("sign failed: %w", err)
	}

	dst := h.Sum(nil)
	return dst, nil
}

func (r *Report) compress(p []byte) (*bytes.Buffer, error) {
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

func (r *Report) isConnRefused(err error) bool {
	var netErr *net.OpError
	return errors.As(err, &netErr) && netErr.Op == "dial" && netErr.Err.Error() == "connect: connection refused"
}

package service

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"

	"github.com/arefev/mtrcstore/internal/agent/model"
	"github.com/arefev/mtrcstore/internal/retry"
	"github.com/go-resty/resty/v2"
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

type Report struct {
	Storage       Storage
	updateURL     string
	massUpdateURL string
	gaugeName     string
	counterName   string
	secretKey     string
	client        resty.Client
}

func NewReport(s Storage, host string, secretKey string) *Report {
	const (
		contentType       = "text/plain"
		protocol          = "http://"
		updateURLPath     = "update"
		massUpdateURLPath = "updates/"
		counterName       = "counter"
		gaugeName         = "gauge"
	)

	client := resty.New().SetBaseURL(protocol + host)
	return &Report{
		Storage:       s,
		updateURL:     updateURLPath,
		massUpdateURL: massUpdateURLPath,
		gaugeName:     gaugeName,
		counterName:   counterName,
		secretKey:     secretKey,
		client:        *client,
	}
}

func (r *Report) Send(metrics []model.Metric) {
	const rCount = 3
	action := func() error {
		return r.request(metrics, r.massUpdateURL)
	}
	if err := retry.New(action, r.isConnRefused, rCount).Run(); err != nil {
		log.Printf("sendCounters(): failed to send the counter metric %s: %s", r.counterName, err.Error())
	}
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

func (r *Report) request(data any, url string) error {
	client := r.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip")

	jsonBody, mErr := json.Marshal(data)
	if mErr != nil {
		return r.requestError(mErr)
	}

	if r.secretKey != "" {
		hash, err := r.sign(jsonBody)
		if err != nil {
			return r.requestError(err)
		}

		client.SetHeader("HashSHA256", hex.EncodeToString(hash))
	}

	body, err := r.compress(jsonBody)
	if err != nil {
		return r.requestError(err)
	}

	if _, err := client.SetBody(body).Post(url); err != nil {
		return r.requestError(err)
	}

	return nil
}

func (r *Report) requestError(err error) error {
	return fmt.Errorf("request failed: %w", err)
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

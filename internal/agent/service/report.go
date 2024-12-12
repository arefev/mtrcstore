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

type report struct {
	Storage       Storage
	updateURL     string
	massUpdateURL string
	gaugeName     string
	counterName   string
	secretKey     string
	client        resty.Client
}

func NewReport(s Storage, host string, secretKey string) (report, error) {
	const (
		contentType       = "text/plain"
		protocol          = "http://"
		updateURLPath     = "update"
		massUpdateURLPath = "updates/"
		counterName       = "counter"
		gaugeName         = "gauge"
	)

	client := resty.New().SetBaseURL(protocol + host)
	return report{
		Storage:       s,
		updateURL:     updateURLPath,
		massUpdateURL: massUpdateURLPath,
		gaugeName:     gaugeName,
		counterName:   counterName,
		secretKey:     secretKey,
		client:        *client,
	}, nil
}

func (r *report) Send() {
	r.sendGauges()
	r.sendCounters()
	r.Storage.ClearCounter()
}

func (r *report) MassSend() error {
	const retryCount = 3
	metrics := r.getMetrics()

	r.Storage.ClearCounter()

	if len(metrics) == 0 {
		return nil
	}

	action := func() error {
		return r.request(metrics, r.massUpdateURL)
	}

	if err := retry.New(action, r.isConnRefused, retryCount).Run(); err != nil {
		return fmt.Errorf("massSend(): failed to send metrics, %w", err)
	}

	return nil
}

func (r *report) getMetrics() []model.Metric {
	metrics := make([]model.Metric, 0)
	metrics = append(metrics, r.getGauges()...)
	metrics = append(metrics, r.getCounters()...)
	return metrics
}

func (r *report) getGauges() []model.Metric {
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

func (r *report) getCounters() []model.Metric {
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

func (r *report) Save(memStats *runtime.MemStats) error {
	if err := r.Storage.Save(memStats); err != nil {
		return fmt.Errorf("report save(): metrics save failed: %w", err)
	}

	return nil
}

func (r *report) SaveCPU() error {
	if err := r.Storage.SaveCPU(); err != nil {
		return fmt.Errorf("report saveCPU() failed: %w", err)
	}

	return nil
}

func (r *report) IncrementCounter() {
	r.Storage.IncrementCounter()
}

func (r *report) sendGauges() {
	for _, metric := range r.getGauges() {
		if err := r.request(metric, r.updateURL); err != nil {
			log.Printf("sendGauges(): failed to send the gauge metric %s: %s", r.gaugeName, err.Error())
			continue
		}
	}
}

func (r *report) sendCounters() {
	for _, metric := range r.getCounters() {
		if err := r.request(metric, r.updateURL); err != nil {
			log.Printf("sendCounters(): failed to send the counter metric %s: %s", r.counterName, err.Error())
			continue
		}
	}
}

func (r *report) request(data any, url string) error {
	client := r.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip")

	jsonBody, err := json.Marshal(data)
	if err != nil {
		return r.requestError(err)
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

func (r *report) requestError(err error) error {
	return fmt.Errorf("request failed: %w", err)
}

func (r *report) sign(data []byte) ([]byte, error) {
	key := []byte(r.secretKey)
	h := hmac.New(sha256.New, key)

	if _, err := h.Write(data); err != nil {
		return []byte{}, fmt.Errorf("sign failed: %w", err)
	}

	dst := h.Sum(nil)
	return dst, nil
}

func (r *report) compress(p []byte) (*bytes.Buffer, error) {
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

func (r *report) isConnRefused(err error) bool {
	var netErr *net.OpError
	return errors.As(err, &netErr) && netErr.Op == "dial" && netErr.Err.Error() == "connect: connection refused"
}

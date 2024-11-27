package service

import (
	"bytes"
	"compress/gzip"
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
	client        resty.Client
}

func NewReport(s Storage, host string) (report, error) {
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
		client:        *client,
	}, nil
}

func (r *report) Send() {
	r.sendGauges()
	r.sendCounters()
	r.Storage.ClearCounter()
}

func (r *report) MassSend() {
	const retryCount = 3
	metrics := make([]model.Metric, 0)
	for name, val := range r.Storage.GetGauges() {
		mVal := float64(val)
		metrics = append(metrics, model.Metric{
			ID:    name,
			MType: r.gaugeName,
			Value: &mVal,
		})
	}

	for name, val := range r.Storage.GetCounters() {
		mVal := int64(val)
		metrics = append(metrics, model.Metric{
			ID:    name,
			MType: r.counterName,
			Delta: &mVal,
		})
	}

	r.Storage.ClearCounter()

	if len(metrics) == 0 {
		return
	}

	action := func() error {
		return r.request(metrics, r.massUpdateURL)
	}

	if err := retry.New(action, r.isConnRefused, retryCount).Run(); err != nil {
		log.Printf("massSend(): failed to send metrics, %s", err.Error())
	}
}

func (r *report) Save(memStats *runtime.MemStats) error {
	if err := r.Storage.Save(memStats); err != nil {
		return fmt.Errorf("report save(): metrics save failed: %w", err)
	}

	return nil
}

func (r *report) IncrementCounter() {
	r.Storage.IncrementCounter()
}

func (r *report) sendGauges() {
	for name, val := range r.Storage.GetGauges() {
		mVal := float64(val)
		metric := model.Metric{
			ID:    name,
			MType: r.gaugeName,
			Value: &mVal,
		}

		if err := r.request(metric, r.updateURL); err != nil {
			log.Printf("sendGauges(): failed to send the gauge metric %s: %s", r.gaugeName, err.Error())
			continue
		}
	}
}

func (r *report) sendCounters() {
	for name, val := range r.Storage.GetCounters() {
		mVal := int64(val)
		metric := model.Metric{
			ID:    name,
			MType: r.counterName,
			Delta: &mVal,
		}

		if err := r.request(metric, r.updateURL); err != nil {
			log.Printf("sendCounters(): failed to send the counter metric %s: %s", r.counterName, err.Error())
			continue
		}
	}
}

func (r *report) request(data any, url string) error {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("send failed: %w", err)
	}

	body, err := r.compress(jsonBody)
	if err != nil {
		return fmt.Errorf("send failed: %w", err)
	}

	_, err = r.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(body).
		Post(url)

	if err != nil {
		return fmt.Errorf("send failed: %w", err)
	}

	return nil
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

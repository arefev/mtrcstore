package service

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"runtime"

	"github.com/arefev/mtrcstore/internal/agent/model"
	"github.com/go-resty/resty/v2"
)

type Gauge float64
type Counter int64

const (
	contentType       = "text/plain"
	protocol          = "http://"
	updateURLPath     = "update"
	massUpdateURLPath = "updates/"
	counterName       = "counter"
	gaugeName         = "gauge"
)

type Storage interface {
	Save(memStats *runtime.MemStats) error
	IncrementCounter()
	ClearCounter()
	GetGauges() map[string]Gauge
	GetCounters() map[string]Counter
}

type Report struct {
	Storage       Storage
	ServerHost    string
	updateURL     string
	massUpdateURL string
	client        resty.Client
}

func NewReport(s Storage, host string) (Report, error) {
	updateURL, err := url.JoinPath(protocol+host, updateURLPath)
	if err != nil {
		return Report{}, fmt.Errorf("NewReport failed: %w", err)
	}

	massUpdateURL, err := url.JoinPath(protocol+host, massUpdateURLPath)
	if err != nil {
		return Report{}, fmt.Errorf("NewReport failed: %w", err)
	}

	client := resty.New()
	return Report{
		Storage:       s,
		ServerHost:    host,
		updateURL:     updateURL,
		massUpdateURL: massUpdateURL,
		client:        *client,
	}, nil
}

func (r *Report) Send() {
	r.sendGauges()
	r.sendCounters()
	r.Storage.ClearCounter()
}

func (r *Report) MassSend() {
	var metrics []model.Metric
	for name, val := range r.Storage.GetGauges() {
		mVal := float64(val)
		metrics = append(metrics, model.Metric{
			ID:    name,
			MType: gaugeName,
			Value: &mVal,
		})
	}

	for name, val := range r.Storage.GetCounters() {
		mVal := int64(val)
		metrics = append(metrics, model.Metric{
			ID:    name,
			MType: counterName,
			Delta: &mVal,
		})
	}

	r.Storage.ClearCounter()
	if err := r.massSend(&metrics); err != nil {
		log.Printf("massSend(): failed to send metrics %+v, %s", metrics, err.Error())
	}
}

func (r *Report) Save(memStats *runtime.MemStats) error {
	if err := r.Storage.Save(memStats); err != nil {
		return fmt.Errorf("report save(): metrics save failed: %w", err)
	}

	return nil
}

func (r *Report) IncrementCounter() {
	r.Storage.IncrementCounter()
}

func (r *Report) sendGauges() {
	for name, val := range r.Storage.GetGauges() {
		mVal := float64(val)
		metric := model.Metric{
			ID:    name,
			MType: gaugeName,
			Value: &mVal,
		}

		if err := r.send(&metric); err != nil {
			log.Printf("sendGauges(): failed to send the gauge metric %s: %s", gaugeName, err.Error())
			continue
		}
	}
}

func (r *Report) sendCounters() {
	for name, val := range r.Storage.GetCounters() {
		mVal := int64(val)
		metric := model.Metric{
			ID:    name,
			MType: counterName,
			Delta: &mVal,
		}

		if err := r.send(&metric); err != nil {
			log.Printf("sendCounters(): failed to send the counter metric %s: %s", counterName, err.Error())
			continue
		}
	}
}

func (r *Report) send(m *model.Metric) error {
	jsonBody, err := json.Marshal(m)
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
		Post(r.updateURL)

	if err != nil {
		return fmt.Errorf("send failed: %w", err)
	}

	return nil
}

func (r *Report) massSend(m *[]model.Metric) error {
	jsonBody, err := json.Marshal(m)
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
		Post(r.massUpdateURL)

	if err != nil {
		return fmt.Errorf("send failed: %w", err)
	}

	return nil
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

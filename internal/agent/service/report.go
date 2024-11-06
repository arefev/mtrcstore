package service

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type Gauge float64
type Counter int64

const (
	contentType   = "text/plain"
	protocol      = "http://"
	updateURLPath = "update"
	counterName   = "counter"
	gaugeName     = "gauge"
)

type Storage interface {
	Save(memStats *runtime.MemStats) error
	IncrementCounter()
	ClearCounter()
	GetGauges() map[string]Gauge
	GetCounters() map[string]Counter
}

type Report struct {
	Storage    Storage
	ServerHost string
}

func NewReport(s Storage, host string) Report {
	return Report{
		Storage:    s,
		ServerHost: host,
	}
}

func (r *Report) Send() {
	r.sendGauges()
	r.sendCounters()
	r.Storage.ClearCounter()
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
	reader := bytes.NewReader([]byte(""))

	for name, val := range r.Storage.GetGauges() {
		url := fmt.Sprintf("%s%s/%s/%s/%s/%f", protocol, r.ServerHost, updateURLPath, gaugeName, name, val)
		resp, err := http.Post(url, contentType, reader)
		if err != nil {
			log.Printf("sendGauges(): failed to send the gauge metric %s: %s", gaugeName, err.Error())
			continue
		}

		if err := resp.Body.Close(); err != nil {
			log.Printf("sendGauges(): body close failed %s: %s", gaugeName, err.Error())
			continue
		}
	}
}

func (r *Report) sendCounters() {
	reader := bytes.NewReader([]byte(""))

	for name, val := range r.Storage.GetCounters() {
		url := fmt.Sprintf("%s%s/%s/%s/%s/%d", protocol, r.ServerHost, updateURLPath, counterName, name, val)
		resp, err := http.Post(url, contentType, reader)
		if err != nil {
			log.Printf("sendCounters(): failed to send the counter metric %s: %s", counterName, err.Error())
			continue
		}

		if err := resp.Body.Close(); err != nil {
			log.Printf("sendCounters(): body close failed %s: %s", counterName, err.Error())
			continue
		}
	}
}

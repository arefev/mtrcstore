package service

import (
	"fmt"
	"log"
	"runtime"

	"github.com/arefev/mtrcstore/internal/agent/model"
	"github.com/go-resty/resty/v2"
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
	client := resty.New()
	url := fmt.Sprintf("%s%s/%s", protocol, r.ServerHost, updateURLPath)

	for name, val := range r.Storage.GetGauges() {
		mVal := float64(val)
		metric := model.Metric{
			ID:    name,
			MType: gaugeName,
			Value: &mVal,
		}
		_, err := client.R().SetBody(metric).Post(url)

		if err != nil {
			log.Printf("sendGauges(): failed to send the gauge metric %s: %s", gaugeName, err.Error())
			continue
		}
	}
}

func (r *Report) sendCounters() {
	client := resty.New()
	url := fmt.Sprintf("%s%s/%s", protocol, r.ServerHost, updateURLPath)

	for name, val := range r.Storage.GetCounters() {
		mVal := int64(val)
		metric := model.Metric{
			ID:    name,
			MType: counterName,
			Delta: &mVal,
		}
		_, err := client.R().SetBody(metric).Post(url)

		if err != nil {
			log.Printf("sendCounters(): failed to send the counter metric %s: %s", counterName, err.Error())
			continue
		}
	}
}

package service

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/arefev/mtrcstore/internal/agent/model"
	"github.com/arefev/mtrcstore/internal/retry"
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
	Request(ctx context.Context, data []model.Metric) error
	IsConnRefused(err error) bool
}

type Report struct {
	Storage     Storage
	sender      Sender
	gaugeName   string
	counterName string
}

func NewReport(s Storage, sender Sender) *Report {
	const (
		counterName = "counter"
		gaugeName   = "gauge"
	)

	return &Report{
		Storage:     s,
		gaugeName:   gaugeName,
		counterName: counterName,
		sender:      sender,
	}
}

func (r *Report) Send(ctx context.Context, metrics []model.Metric) {
	const rCount = 3
	action := func() error {
		return r.sender.Request(ctx, metrics)
	}
	if err := retry.New(action, r.sender.IsConnRefused, rCount).Run(); err != nil {
		log.Printf("report failed to send the metrics: %s", err.Error())
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

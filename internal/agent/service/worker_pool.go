package service

import (
	"github.com/arefev/mtrcstore/internal/agent/model"
)

type WorkerPool struct {
	Report    *Report
	jobChan   chan []model.Metric
	rateLimit int
}

func NewWorkerPool(report *Report, rateLimit int) *WorkerPool {
	return &WorkerPool{
		Report:    report,
		rateLimit: rateLimit,
	}
}

func (wp *WorkerPool) Run() {
	wp.jobChan = make(chan []model.Metric, wp.rateLimit)

	for range wp.rateLimit {
		go wp.worker()
	}
}

func (wp *WorkerPool) worker() {
	for metrics := range wp.jobChan {
		wp.Report.Send(metrics)
	}
}

func (wp *WorkerPool) Send() {
	wp.jobChan <- wp.Report.GetMetrics()
	wp.Report.ClearCounter()
}

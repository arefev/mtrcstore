package service

import (
	"context"

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

func (wp *WorkerPool) Run(ctx context.Context) {
	wp.jobChan = make(chan []model.Metric, wp.rateLimit)

	for range wp.rateLimit {
		go wp.worker(ctx)
	}
}

func (wp *WorkerPool) worker(ctx context.Context) {
	for metrics := range wp.jobChan {
		wp.Report.Send(ctx, metrics)
	}
}

func (wp *WorkerPool) Send() {
	wp.jobChan <- wp.Report.GetMetrics()
	wp.Report.ClearCounter()
}

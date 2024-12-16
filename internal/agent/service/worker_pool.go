package service

import (
	"github.com/arefev/mtrcstore/internal/agent/model"
)

type WorkerPool struct {
	Report    *Report
	jobChan   chan model.Metric
	rateLimit int
}

func NewWorkerPool(report *Report, rateLimit int) *WorkerPool {
	return &WorkerPool{
		Report:    report,
		rateLimit: rateLimit,
	}
}

func (wp *WorkerPool) Run() {
	wp.jobChan = make(chan model.Metric, wp.rateLimit)

	for range wp.rateLimit {
		go wp.worker()
	}
}

func (wp *WorkerPool) worker() {
	for metric := range wp.jobChan {
		wp.Report.Send(metric)
	}
}

func (wp *WorkerPool) Send() {
	for _, m := range wp.Report.GetMetrics() {
		wp.jobChan <- m
	}

	wp.Report.ClearCounter()
}

func (wp *WorkerPool) IsEmpty() bool {
	return len(wp.jobChan) == 0
}

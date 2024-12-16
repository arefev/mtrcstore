package agent

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/arefev/mtrcstore/internal/agent/service"
	"golang.org/x/sync/errgroup"
)

type Worker struct {
	WorkerPool     *service.WorkerPool
	PollInterval   int
	ReportInterval int
}

func (w *Worker) Run() error {
	var memStats runtime.MemStats
	readTime := time.NewTicker(time.Duration(w.PollInterval) * time.Second).C
	sendTime := time.NewTicker(time.Duration(w.ReportInterval) * time.Second).C

	w.WorkerPool.Run()

	for {
		select {
		case <-readTime:
			log.Println("readTime")
			if err := w.read(&memStats); err != nil {
				return fmt.Errorf("Worker Run() failed: %w", err)
			}
		case <-sendTime:
			log.Println("sendTime")
			w.WorkerPool.Send()
		}
	}
}

func (w *Worker) read(memStats *runtime.MemStats) error {
	g := &errgroup.Group{}

	g.Go(func() error {
		w.WorkerPool.Report.IncrementCounter()
		runtime.ReadMemStats(memStats)
		if err := w.WorkerPool.Report.Save(memStats); err != nil {
			return fmt.Errorf("worker read(): metrics save failed: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := w.WorkerPool.Report.SaveCPU(); err != nil {
			return fmt.Errorf("worker read(): CPU save failed: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("worker read(): read metrics failed: %w", err)
	}

	return nil
}

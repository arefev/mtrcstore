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
	var period int
	var memStats runtime.MemStats
	start := time.Now()

	w.WorkerPool.Run()

	for {
		if err := w.read(&memStats); err != nil {
			return fmt.Errorf("Worker Run() failed: %w", err)
		}

		time.Sleep(time.Duration(w.PollInterval * int(time.Second)))

		period = int(time.Until(start).Abs().Seconds())

		if period >= w.ReportInterval {
			log.Printf("Run report after %d seconds", period)

			w.WorkerPool.Send()
			start = time.Now()
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

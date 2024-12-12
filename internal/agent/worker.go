package agent

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Reporter interface {
	Send()
	MassSend() error
	PoolSend()
	Save(memStats *runtime.MemStats) error
	SaveCPU() error
	IncrementCounter()
}

type Worker struct {
	Report         Reporter
	PollInterval   int
	ReportInterval int
}

func (w *Worker) Run() error {
	var period int
	var memStats runtime.MemStats
	start := time.Now()

	for {
		if err := w.read(&memStats); err != nil {
			return fmt.Errorf("Worker Run() failed: %w", err)
		}

		time.Sleep(time.Duration(w.PollInterval * int(time.Second)))

		period = int(time.Until(start).Abs().Seconds())

		if period >= w.ReportInterval {
			log.Printf("Run report after %d seconds", period)

			rTimeStart := time.Now()

			w.Report.PoolSend()
			// w.Report.Send()

			rDuration := time.Since(rTimeStart)
			log.Printf("request duration: %v", rDuration)

			// if err := w.Report.MassSend(); err != nil {
			// 	log.Printf("worker failed: %v", err)
			// }
			start = time.Now()
		}
	}
}

func (w *Worker) read(memStats *runtime.MemStats) error {
	var m sync.Mutex
	g := &errgroup.Group{}

	g.Go(func() error {
		m.Lock()
		defer m.Unlock()
		w.Report.IncrementCounter()
		runtime.ReadMemStats(memStats)
		if err := w.Report.Save(memStats); err != nil {
			return fmt.Errorf("worker read(): metrics save failed: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		m.Lock()
		defer m.Unlock()
		if err := w.Report.SaveCPU(); err != nil {
			return fmt.Errorf("worker read(): CPU save failed: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("worker read(): read metrics failed: %w", err)
	}

	return nil
}

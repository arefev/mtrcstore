package agent

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

type Reporter interface {
	Send()
	Save(memStats *runtime.MemStats) error
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
		w.Report.IncrementCounter()
		if err := w.read(&memStats); err != nil {
			return fmt.Errorf("Worker Run() failed: %w", err)
		}

		time.Sleep(time.Duration(w.PollInterval * int(time.Second)))

		period = int(time.Until(start).Abs().Seconds())

		if period >= w.ReportInterval {
			log.Printf("Run report after %d seconds", period)

			w.Report.Send()
			start = time.Now()
		}
	}
}

func (w *Worker) read(memStats *runtime.MemStats) error {
	runtime.ReadMemStats(memStats)
	if err := w.Report.Save(memStats); err != nil {
		return fmt.Errorf("worker read(): metrics save failed: %w", err)
	}

	return nil
}

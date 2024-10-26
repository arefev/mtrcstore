package agent

import (
	"fmt"
	"runtime"
	"time"

	"github.com/arefev/mtrcstore/internal/agent/repository"
)

type Worker struct {
	ReportInterval float64
	PollInterval int
	Storage repository.Storage
}

func (w *Worker) Run() {
	var period float64
	var memStats runtime.MemStats
	start := time.Now()
	
	for {
		w.Storage.IncrementCounter()

		w.read(&memStats)

		time.Sleep(time.Duration(w.PollInterval * int(time.Second)))

		period = time.Until(start).Abs().Seconds()

		fmt.Printf("\nDuration %f\n", period)
		if period >= w.ReportInterval {
			w.report()
			start = time.Now()
			w.Storage.ClearCounter()
		}
	}
}

func (w *Worker) read(memStats *runtime.MemStats) {
	runtime.ReadMemStats(memStats)

	w.Storage.Save(memStats)

	fmt.Printf("%+v\n", w.Storage)
}

func (w *Worker) report() {
	fmt.Println("\nSend metrics to server")
	fmt.Printf("Data = %v\n\n", w.Storage)
}
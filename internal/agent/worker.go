package agent

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/arefev/mtrcstore/internal/agent/repository"
)

const contentType = "text/plain"
const protocol = "http://"

type Worker struct {
	ReportInterval int
	PollInterval   int
	Storage        repository.Storage
	ServerHost     string
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

		if int(period) >= w.ReportInterval {
			w.report()
			start = time.Now()
			w.Storage.ClearCounter()
		}
	}
}

func (w *Worker) read(memStats *runtime.MemStats) {
	runtime.ReadMemStats(memStats)
	w.Storage.Save(memStats)
}

func (w *Worker) report() {
	w.sendGauges()
	w.sendCounters()
}

func (w *Worker) getReportUrl(mType string, name string, val float64) string {
	return fmt.Sprintf("%s%s/update/%s/%s/%f", protocol, w.ServerHost, mType, name, val)
}

func (w *Worker) sendGauges() {
	const mType = "gauge"
	r := bytes.NewReader([]byte(""))

	for name, val := range w.Storage.GetGauges() {
		qPath := w.getReportUrl(mType, name, float64(val))
		if _, err := http.Post(qPath, contentType, r); err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func (w *Worker) sendCounters() {
	const mType = "counter"
	r := bytes.NewReader([]byte(""))

	for name, val := range w.Storage.GetCounters() {
		qPath := w.getReportUrl(mType, name, float64(val))
		if _, err := http.Post(qPath, contentType, r); err != nil {
			fmt.Println(err)
			continue
		}
	}
}

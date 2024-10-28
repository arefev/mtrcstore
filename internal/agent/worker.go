package agent

import (
	"bytes"
	"fmt"
	"log"
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
	var period int
	var memStats runtime.MemStats
	start := time.Now()

	for {
		w.Storage.IncrementCounter()
		w.read(&memStats)

		time.Sleep(time.Duration(w.PollInterval * int(time.Second)))

		period = int(time.Until(start).Abs().Seconds())

		if period >= w.ReportInterval {
			log.Printf("Run report after %d seconds", period)

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

func (w *Worker) getReportURL(mType string, name string, val float64) string {
	return fmt.Sprintf("%s%s/update/%s/%s/%f", protocol, w.ServerHost, mType, name, val)
}

func (w *Worker) sendGauges() {
	const mType = "gauge"
	r := bytes.NewReader([]byte(""))

	for name, val := range w.Storage.GetGauges() {
		qPath := w.getReportURL(mType, name, float64(val))
		resp, err := http.Post(qPath, contentType, r)
		if err != nil {
			log.Fatal(err)
			continue
		}
		resp.Body.Close()
	}
}

func (w *Worker) sendCounters() {
	const mType = "counter"
	r := bytes.NewReader([]byte(""))

	for name, val := range w.Storage.GetCounters() {
		qPath := w.getReportURL(mType, name, float64(val))
		resp, err := http.Post(qPath, contentType, r)
		if err != nil {
			log.Fatal(err)
			continue
		}
		resp.Body.Close()
	}
}

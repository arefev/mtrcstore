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

const (
	contentType = "text/plain"
	protocol = "http://"
	updateUrlPath = "update"
	counterName = "counter"
	gaugeName = "gauge"
)

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

func (w *Worker) sendGauges() {
	r := bytes.NewReader([]byte(""))

	for name, val := range w.Storage.GetGauges() {
		url := fmt.Sprintf("%s%s/%s/%s/%s/%f", protocol, w.ServerHost, updateUrlPath, gaugeName, name, val)
		resp, err := http.Post(url, contentType, r)
		if err != nil {
			log.Print(err)
			continue
		}
		resp.Body.Close()
	}
}

func (w *Worker) sendCounters() {
	r := bytes.NewReader([]byte(""))

	for name, val := range w.Storage.GetCounters() {
		url := fmt.Sprintf("%s%s/%s/%s/%s/%d", protocol, w.ServerHost, updateUrlPath, counterName, name, val)
		resp, err := http.Post(url, contentType, r)
		if err != nil {
			log.Print(err)
			continue
		}
		resp.Body.Close()
	}
}

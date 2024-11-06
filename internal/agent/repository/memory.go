package repository

import (
	"math/rand"
	"runtime"

	"github.com/arefev/mtrcstore/internal/agent/service"
)

type memory struct {
	Gauge   map[string]service.Gauge
	Counter map[string]service.Counter
}

func NewMemory() memory {
	return memory{
		Gauge:   make(map[string]service.Gauge),
		Counter: make(map[string]service.Counter),
	}
}

func (s *memory) Save(memStats *runtime.MemStats) error {
	s.Gauge["Alloc"] = service.Gauge(memStats.Alloc)
	s.Gauge["BuckHashSys"] = service.Gauge(memStats.BuckHashSys)
	s.Gauge["Frees"] = service.Gauge(memStats.Frees)
	s.Gauge["GCCPUFraction"] = service.Gauge(memStats.GCCPUFraction)
	s.Gauge["GCSys"] = service.Gauge(memStats.GCSys)
	s.Gauge["HeapAlloc"] = service.Gauge(memStats.HeapAlloc)
	s.Gauge["HeapIdle"] = service.Gauge(memStats.HeapIdle)
	s.Gauge["HeapInuse"] = service.Gauge(memStats.HeapInuse)
	s.Gauge["HeapObjects"] = service.Gauge(memStats.HeapObjects)
	s.Gauge["HeapReleased"] = service.Gauge(memStats.HeapReleased)
	s.Gauge["HeapSys"] = service.Gauge(memStats.HeapSys)
	s.Gauge["LastGC"] = service.Gauge(memStats.LastGC)
	s.Gauge["Lookups"] = service.Gauge(memStats.Lookups)
	s.Gauge["MCacheInuse"] = service.Gauge(memStats.MCacheInuse)
	s.Gauge["MCacheSys"] = service.Gauge(memStats.MCacheSys)
	s.Gauge["MSpanInuse"] = service.Gauge(memStats.MSpanInuse)
	s.Gauge["MSpanSys"] = service.Gauge(memStats.MSpanSys)
	s.Gauge["Mallocs"] = service.Gauge(memStats.Mallocs)
	s.Gauge["NextGC"] = service.Gauge(memStats.NextGC)
	s.Gauge["NumForcedGC"] = service.Gauge(memStats.NumForcedGC)
	s.Gauge["NumGC"] = service.Gauge(memStats.NumGC)
	s.Gauge["OtherSys"] = service.Gauge(memStats.OtherSys)
	s.Gauge["PauseTotalNs"] = service.Gauge(memStats.PauseTotalNs)
	s.Gauge["StackInuse"] = service.Gauge(memStats.StackInuse)
	s.Gauge["StackSys"] = service.Gauge(memStats.StackSys)
	s.Gauge["Sys"] = service.Gauge(memStats.Sys)
	s.Gauge["TotalAlloc"] = service.Gauge(memStats.TotalAlloc)
	s.Gauge["RandomValues"] = service.Gauge(rand.Int())
	return nil
}

func (s *memory) IncrementCounter() {
	s.Counter["PollCount"]++
}

func (s *memory) ClearCounter() {
	s.Counter["PollCount"] = 0
}

func (s *memory) GetGauges() map[string]service.Gauge {
	return s.Gauge
}

func (s *memory) GetCounters() map[string]service.Counter {
	return s.Counter
}

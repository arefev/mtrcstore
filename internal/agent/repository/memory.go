package repository

import (
	"math/rand"
	"runtime"
)

type gauge float64
type counter int64

type memory struct {
	Gauge   map[string]gauge
	Counter map[string]counter
}

func NewMemory() memory {
	return memory{
		Gauge:   make(map[string]gauge),
		Counter: make(map[string]counter),
	}
}

func (s *memory) Save(memStats *runtime.MemStats) error {
	s.Gauge["Alloc"] = gauge(memStats.Alloc)
	s.Gauge["BuckHashSys"] = gauge(memStats.BuckHashSys)
	s.Gauge["Frees"] = gauge(memStats.Frees)
	s.Gauge["GCCPUFraction"] = gauge(memStats.GCCPUFraction)
	s.Gauge["GCSys"] = gauge(memStats.GCSys)
	s.Gauge["HeapAlloc"] = gauge(memStats.HeapAlloc)
	s.Gauge["HeapIdle"] = gauge(memStats.HeapIdle)
	s.Gauge["HeapInuse"] = gauge(memStats.HeapInuse)
	s.Gauge["HeapObjects"] = gauge(memStats.HeapObjects)
	s.Gauge["HeapReleased"] = gauge(memStats.HeapReleased)
	s.Gauge["HeapSys"] = gauge(memStats.HeapSys)
	s.Gauge["LastGC"] = gauge(memStats.LastGC)
	s.Gauge["Lookups"] = gauge(memStats.Lookups)
	s.Gauge["MCacheInuse"] = gauge(memStats.MCacheInuse)
	s.Gauge["MCacheSys"] = gauge(memStats.MCacheSys)
	s.Gauge["MSpanInuse"] = gauge(memStats.MSpanInuse)
	s.Gauge["MSpanSys"] = gauge(memStats.MSpanSys)
	s.Gauge["Mallocs"] = gauge(memStats.Mallocs)
	s.Gauge["NextGC"] = gauge(memStats.NextGC)
	s.Gauge["NumForcedGC"] = gauge(memStats.NumForcedGC)
	s.Gauge["NumGC"] = gauge(memStats.NumGC)
	s.Gauge["OtherSys"] = gauge(memStats.OtherSys)
	s.Gauge["PauseTotalNs"] = gauge(memStats.PauseTotalNs)
	s.Gauge["StackInuse"] = gauge(memStats.StackInuse)
	s.Gauge["StackSys"] = gauge(memStats.StackSys)
	s.Gauge["Sys"] = gauge(memStats.Sys)
	s.Gauge["TotalAlloc"] = gauge(memStats.TotalAlloc)
	s.Gauge["RandomValue"] = gauge(rand.Int())
	return nil
}

func (s *memory) IncrementCounter() {
	s.Counter["PollCount"]++
}

func (s *memory) ClearCounter() {
	s.Counter["PollCount"] = 0
}

func (s *memory) GetGauges() map[string]gauge {
	return s.Gauge
}

func (s *memory) GetCounters() map[string]counter {
	return s.Counter
}

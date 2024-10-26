package repository

import "runtime"

type gauge float64
type counter int64

type Memory struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	PollCount     counter
	RandomValue   counter
}

func (s *Memory) Save(memStats *runtime.MemStats) error {
	s.Alloc = gauge(memStats.Alloc)
	return nil
}

func (s *Memory) IncrementCounter() {
	s.PollCount++
}

func (s *Memory) ClearCounter() {
	s.PollCount = 0
}

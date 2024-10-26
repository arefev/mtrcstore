package repository

import "runtime"

type Storage interface {
	Save(memStats *runtime.MemStats) error
	IncrementCounter()
	ClearCounter()
	GetGauges() map[string]gauge
	GetCounters() map[string]counter
}

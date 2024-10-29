package repository

import "errors"

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

func (s *memory) Save(mType string, name string, value float64) error {

	switch mType {
	case "counter":
		s.Counter[name] += counter(value)
	default:
		s.Gauge[name] = gauge(value)
	}

	return nil
}

func (s *memory) Find(mType string, name string) (float64, error) {
	switch mType {
	case "counter":
		val, ok := s.Counter[name]
		if !ok {
			return 0, errors.New("counter value not found")
		}

		return float64(val), nil
	default:
		val, ok := s.Gauge[name]
		if !ok {
			return 0, errors.New("gauge value not found")
		}

		return float64(val), nil
	}
}

func (s *memory) Get() map[string]float64 {
	all := make(map[string]float64)
	for name, val := range s.Gauge {
		all[name] = float64(val)
	}

	for name, val := range s.Counter {
		all[name] = float64(val)
	}

	return all
}

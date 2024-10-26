package repository

import "strconv"

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

func (s *memory) Save(mType string, name string, value string) error {
	
	switch mType {
	case "counter":
		val, err := strconv.Atoi(value)
		if err != nil {
			
		}
		s.Counter[name] = counter(val)
	default:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {

		}
		s.Gauge[name] = gauge(val)
	}

	return nil
}

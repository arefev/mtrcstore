package repository

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

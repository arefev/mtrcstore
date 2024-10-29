package repository

import (
	"errors"
	"strconv"
)

type gauge float64
type counter int64

func (g gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', 3, 64)
}

func (c counter) String() string {
	return strconv.Itoa(int(c))
}

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

func (s *memory) FindGauge(name string) (gauge, error) {
	val, ok := s.Gauge[name]
	if !ok {
		return 0, errors.New("gauge value not found")
	}

	return val, nil
}

func (s *memory) FindCounter(name string) (counter, error) {
	val, ok := s.Counter[name]
	if !ok {
		return 0, errors.New("counter value not found")
	}

	return val, nil
}

func (s *memory) Get() map[string]string {
	all := make(map[string]string)
	for name, val := range s.Gauge {
		all[name] = val.String()
	}

	for name, val := range s.Counter {
		all[name] = val.String()
	}

	return all
}

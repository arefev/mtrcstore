package repository

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/arefev/mtrcstore/internal/server/model"
)

const (
	CounterName = "counter"
	GaugeName   = "gauge"
)

type gauge float64
type counter int64

func (g gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
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

func (s *memory) Save(m model.Metric) error {
	switch m.MType {
	case "counter":
		if m.Delta == nil {
			return errors.New("counter has not value")
		}
		s.Counter[m.ID] += counter(*m.Delta)
	default:
		if m.Value == nil {
			return errors.New("gauge has not value")
		}

		s.Gauge[m.ID] = gauge(*m.Value)
	}

	return nil
}

func (s *memory) FindGauge(name string) (model.Metric, error) {
	val, ok := s.Gauge[name]
	if !ok {
		return model.Metric{}, fmt.Errorf("gauge with name %s not found", name)
	}

	value := float64(val)
	metric := model.Metric{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}

	return metric, nil
}

func (s *memory) FindCounter(name string) (model.Metric, error) {
	val, ok := s.Counter[name]
	if !ok {
		return model.Metric{}, fmt.Errorf("counter with name %s not found", name)
	}

	value := int64(val)
	metric := model.Metric{
		ID:    name,
		MType: "counter",
		Delta: &value,
	}

	return metric, nil
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

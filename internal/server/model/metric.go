package model

import "strconv"

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metric) ValueString() string {
	return strconv.FormatFloat(float64(*m.Value), 'f', -1, 64)
}

func (m *Metric) DeltaString() string {
	return strconv.Itoa(int(*m.Delta))
}

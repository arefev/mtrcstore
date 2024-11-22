package model

import "strconv"

type Metric struct {
	Delta *int64   `json:"delta,omitempty" db:"delta"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty" db:"value"` // значение метрики в случае передачи gauge
	ID    string   `json:"id" db:"name"`               // имя метрики
	MType string   `json:"type" db:"type"`             // параметр, принимающий значение gauge или counter
}

func (m *Metric) ValueString() string {
	return strconv.FormatFloat(float64(*m.Value), 'f', -1, 64)
}

func (m *Metric) DeltaString() string {
	return strconv.Itoa(int(*m.Delta))
}

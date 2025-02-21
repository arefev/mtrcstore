package model

import "strconv"

type Metric struct {
	Delta *int64   `json:"delta,omitempty" db:"delta"` // metric value in case of counter transfer
	Value *float64 `json:"value,omitempty" db:"value"` // metric value in case of gauge transfer
	ID    string   `json:"id" db:"name"`               // metric name
	MType string   `json:"type" db:"type"`             // parameter that takes the value gauge or counter
}

func (m *Metric) ValueString() string {
	return strconv.FormatFloat(float64(*m.Value), 'f', -1, 64)
}

func (m *Metric) DeltaString() string {
	return strconv.Itoa(int(*m.Delta))
}

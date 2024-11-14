package model

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Metric struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
}

func (m *Metric) ToJSON() (string, error) {
	b := bytes.NewBuffer(nil)
	decoder := json.NewEncoder(b)
	if err := decoder.Encode(m); err != nil {
		return "", fmt.Errorf("model Metric json encode failed: %w", err)
	}

	return b.String(), nil
}
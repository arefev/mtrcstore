package repository

import "github.com/arefev/mtrcstore/internal/server/model"

type Storage interface {
	Save(m model.Metric) error
	FindGauge(name string) (model.Metric, error)
	FindCounter(name string) (model.Metric, error)
	Get() map[string]string
}

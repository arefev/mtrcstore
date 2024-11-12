package repository

import "github.com/arefev/mtrcstore/internal/server/model"

type Storage interface {
	Save(m model.Metric) error
	FindGauge(name string) (gauge, error)
	FindCounter(name string) (counter, error)
	Get() map[string]string
}

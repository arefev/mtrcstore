package repository

import "github.com/arefev/mtrcstore/internal/server/model"

type Storage interface {
	Save(m model.Metric) error
	Find(id string, mType string) (model.Metric, error)
	Get() map[string]string
	Ping() error
}

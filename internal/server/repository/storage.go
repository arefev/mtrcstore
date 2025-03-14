package repository

import (
	"context"

	"github.com/arefev/mtrcstore/internal/server/model"
)

type Storage interface {
	Save(ctx context.Context, m model.Metric) error
	MassSave(ctx context.Context, elems []model.Metric) error
	Find(ctx context.Context, id string, mType string) (model.Metric, error)
	Get(ctx context.Context) map[string]string
	Ping(ctx context.Context) error
	Close() error
}

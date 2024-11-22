package repository

import (
	"context"
	"database/sql"

	"github.com/arefev/mtrcstore/internal/server/model"
)

type databaseRep struct {
	db *sql.DB
}

func NewDatabaseRep(db *sql.DB) *databaseRep {
	return &databaseRep{
		db: db,
	}
}

func (rep *databaseRep) Save(m model.Metric) error {
	return nil
}

func (rep *databaseRep) Find(id string, mType string) (model.Metric, error) {
	return model.Metric{}, nil
}

func (rep *databaseRep) Get() map[string]string {
	return map[string]string{}
}

func (rep *databaseRep) Ping() error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	return rep.db.PingContext(ctx)
}

package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/arefev/mtrcstore/internal/server/model"
)

type databaseRep struct {
	db *sql.DB
}

func NewDatabaseRep(dsn string) (*databaseRep, error) {
	rep := &databaseRep{}
	if err := rep.connect(dsn); err != nil {
		return &databaseRep{}, err
	}

	return rep, nil
}

func (rep *databaseRep) connect(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("db init failed: %w", err)
	}

	rep.db = db
	return nil
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

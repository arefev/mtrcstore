package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

	if err := rep.migrations(); err != nil {
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

func (rep *databaseRep) migrations() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 1 * time.Second)
	defer cancel()
	
	if err := rep.createTableMetrics(ctx); err != nil {
		return err
	}

	return nil
}

func (rep *databaseRep) createTableMetrics(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS public.metrics (
			id bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
			"type" varchar NULL,
			"name" varchar NULL,
			value double precision NULL,
			delta int NULL,
			CONSTRAINT metrics_pk PRIMARY KEY (id)
		);
	`

	_, err := rep.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("create table metrics failed: %w", err)
	}

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

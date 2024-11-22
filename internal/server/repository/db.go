package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"errors"

	"github.com/arefev/mtrcstore/internal/server/model"
	"go.uber.org/zap"
)

type databaseRep struct {
	db  *sql.DB
	log *zap.Logger
}

func NewDatabaseRep(dsn string, log *zap.Logger) (*databaseRep, error) {
	rep := &databaseRep{
		log: log,
	}

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
	const timeCancel = time.Duration(1)
	ctx, cancel := context.WithTimeout(context.TODO(), timeCancel*time.Second)
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
		return fmt.Errorf("rep db createTableMetrics failed: %w", err)
	}

	return nil
}

func (rep *databaseRep) Save(m model.Metric) error {
	const timeCancel = time.Duration(1)
	ctx, cancel := context.WithTimeout(context.TODO(), timeCancel*time.Second)
	defer cancel()

	metric, err := rep.Find(m.ID, m.MType)
	notFound := errors.Is(err, sql.ErrNoRows)

	if notFound {
		return rep.create(ctx, m)
	} else if err == nil {
		return rep.update(ctx, m, metric)
	}

	return fmt.Errorf("rep db Save failed: %w", err)
}

func (rep *databaseRep) create(ctx context.Context, m model.Metric) error {
	query := "INSERT INTO metrics(type, name, value, delta) VALUES($1, $2, $3, $4)"
	_, err := rep.db.ExecContext(ctx, query, m.MType, m.ID, m.Value, m.Delta)
	if err != nil {
		return fmt.Errorf("rep db create failed: %w", err)
	}

	return nil
}

func (rep *databaseRep) update(ctx context.Context, newMetric model.Metric, oldMetric model.Metric) error {
	query := "UPDATE metrics SET value = $1, delta = $2 WHERE type = $3 AND name = $4"
	_, err := rep.db.ExecContext(ctx, query, newMetric.Value, newMetric.Delta, oldMetric.MType, oldMetric.ID)
	if err != nil {
		return fmt.Errorf("rep db update failed: %w", err)
	}

	return nil
}

func (rep *databaseRep) Find(id string, mType string) (model.Metric, error) {
	const timeCancel = time.Duration(1)
	metric := model.Metric{}
	ctx, cancel := context.WithTimeout(context.TODO(), timeCancel*time.Second)
	defer cancel()

	query := "SELECT type, name, value, delta FROM metrics WHERE type = $1 AND name = $2"
	row := rep.db.QueryRowContext(ctx, query, mType, id)

	err := row.Err()
	if err != nil {
		rep.log.Error("rep db Find failed", zap.Error(err))
		return model.Metric{}, fmt.Errorf("rep db Find failed: %w", err)
	}

	err = row.Scan(&metric.MType, &metric.ID, &metric.Value, &metric.Delta)
	if err != nil {
		rep.log.Error("rep db Find failed", zap.Error(err))
		return model.Metric{}, fmt.Errorf("rep db Find failed: %w", err)
	}

	return metric, nil
}

func (rep *databaseRep) Get() map[string]string {
	const timeCancel = time.Duration(1)
	list := make(map[string]string)
	ctx, cancel := context.WithTimeout(context.TODO(), timeCancel*time.Second)
	defer cancel()

	query := "SELECT type, name, value, delta FROM metrics ORDER BY type, name ASC"
	rows, err := rep.db.QueryContext(ctx, query)
	if err != nil {
		rep.log.Error("rep db Get failed", zap.Error(err))
		return map[string]string{}
	}
	defer rows.Close()

	for rows.Next() {
		var m model.Metric
		err := rows.Scan(&m.MType, &m.ID, &m.Value, &m.Delta)
		if err != nil {
			rep.log.Error("rep db Get failed", zap.Error(err))
			return map[string]string{}
		}

		switch m.MType {
		case "counter":
			list[m.ID] = m.DeltaString()
		default:
			list[m.ID] = m.ValueString()
		}
	}

	if err := rows.Err(); err != nil {
		rep.log.Error("rep db Get failed", zap.Error(err))
		return map[string]string{}
	}

	return list
}

func (rep *databaseRep) Ping() error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	return rep.db.PingContext(ctx)
}

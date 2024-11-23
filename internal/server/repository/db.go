package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/arefev/mtrcstore/internal/server/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type databaseRep struct {
	db  *sqlx.DB
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
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return fmt.Errorf("db init failed: %w", err)
	}

	rep.db = db
	return nil
}

func (rep *databaseRep) migrations() error {
	const timeCancel = 1 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), timeCancel)
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
			delta bigint NULL,
			CONSTRAINT metrics_pk PRIMARY KEY (id),
			CONSTRAINT metrics_unique UNIQUE (type, name)
		);
	`

	_, err := rep.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("rep db createTableMetrics failed: %w", err)
	}

	return nil
}

func (rep *databaseRep) Save(m model.Metric) error {
	const timeCancel = 1 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), timeCancel)
	defer cancel()

	metric, err := rep.Find(m.ID, m.MType)

	switch {
	case err == nil:
		return rep.update(ctx, m, metric)
	case errors.Is(err, sql.ErrNoRows):
		return rep.create(ctx, m)
	default:
		return fmt.Errorf("rep db Save failed: %w", err)
	}
}

func (rep *databaseRep) MassSave(elems []model.Metric) error {
	const timeCancel = 1 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), timeCancel)
	defer cancel()

	tx := rep.db.MustBegin()
	defer func() {
		if err := tx.Rollback(); err != nil {
			rep.log.Warn("rep db mass save rollback failed", zap.Error(err))
		}
	}()

	query := `
		INSERT INTO 
			metrics (type, name, value, delta) 
		VALUES (:type, :name, :value, :delta) 
		ON CONFLICT (type, name)
		DO UPDATE 
		SET value = EXCLUDED.value, delta = EXCLUDED.delta + metrics.delta
	`

	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("rep db mass save failed: %w", err)
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			rep.log.Warn("rep db mass save failed", zap.Error(err))
		}
	}()

	for _, m := range elems {
		_, err := stmt.ExecContext(
			ctx,
			map[string]interface{}{
				"type":  m.MType,
				"name":  m.ID,
				"value": m.Value,
				"delta": m.Delta,
			},
		)

		if err != nil {
			rep.log.Error("rep db mass save failed", zap.Error(err))
			return fmt.Errorf("rep db mass save failed: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		rep.log.Error("rep db mass save commit failed", zap.Error(err))
		return fmt.Errorf("rep db mass save commit failed: %w", err)
	}

	return nil
}

func (rep *databaseRep) create(ctx context.Context, m model.Metric) error {
	query := "INSERT INTO metrics(type, name, value, delta) VALUES(:type, :name, :value, :delta)"

	_, err := rep.db.NamedExecContext(
		ctx,
		query,
		map[string]interface{}{
			"type":  m.MType,
			"name":  m.ID,
			"value": m.Value,
			"delta": m.Delta,
		},
	)
	if err != nil {
		rep.log.Error("rep db create failed", zap.Error(err))
		return fmt.Errorf("rep db create failed: %w", err)
	}

	return nil
}

func (rep *databaseRep) update(ctx context.Context, newMetric model.Metric, oldMetric model.Metric) error {
	query := "UPDATE metrics SET value = :value, delta = :delta WHERE type = :type AND name = :name"
	if oldMetric.MType == "counter" {
		newVal := *oldMetric.Delta + *newMetric.Delta
		newMetric.Delta = &newVal
	}

	_, err := rep.db.NamedExecContext(
		ctx,
		query,
		map[string]interface{}{
			"type":  oldMetric.MType,
			"name":  oldMetric.ID,
			"value": newMetric.Value,
			"delta": newMetric.Delta,
		},
	)

	if err != nil {
		return fmt.Errorf("rep db update failed: %w", err)
	}

	return nil
}

func (rep *databaseRep) Find(id string, mType string) (model.Metric, error) {
	const timeCancel = 1 * time.Second
	metric := model.Metric{}
	ctx, cancel := context.WithTimeout(context.TODO(), timeCancel)
	defer cancel()

	query := "SELECT type, name, value, delta FROM metrics WHERE type = $1 AND name = $2"
	err := rep.db.GetContext(ctx, &metric, query, mType, id)

	if err != nil {
		return model.Metric{}, fmt.Errorf("rep db Find failed: %w", err)
	}

	return metric, nil
}

func (rep *databaseRep) Get() map[string]string {
	const timeCancel = 1 * time.Second
	list := make(map[string]string)
	ctx, cancel := context.WithTimeout(context.TODO(), timeCancel)
	defer cancel()

	query := "SELECT type, name, value, delta FROM metrics ORDER BY type, name ASC"
	metrics := []model.Metric{}
	err := rep.db.SelectContext(ctx, &metrics, query)

	if err != nil {
		rep.log.Error("rep db Get failed", zap.Error(err))
		return map[string]string{}
	}

	for _, m := range metrics {
		switch m.MType {
		case "counter":
			list[m.ID] = m.DeltaString()
		default:
			list[m.ID] = m.ValueString()
		}
	}

	return list
}

func (rep *databaseRep) Ping() error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	if err := rep.db.PingContext(ctx); err != nil {
		return fmt.Errorf("Ping failed: %w", err)
	}

	return nil
}

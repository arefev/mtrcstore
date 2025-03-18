package repository

import (
	"context"
	"testing"

	"github.com/arefev/mtrcstore/internal/server/logger"
	"github.com/arefev/mtrcstore/internal/server/model"
	"github.com/arefev/mtrcstore/internal/server/repository/testdb"
	"github.com/stretchr/testify/require"
)

func TestDBClose(t *testing.T) {
	t.Run("db close success", func(t *testing.T) {
		ctx := context.Background()

		testDB, err := testdb.New(ctx)
		require.NoError(t, err)

		cLog, err := logger.Build("debug")
		require.NoError(t, err)

		rep, err := NewDatabaseRep(testDB.URI, cLog)
		require.NoError(t, err)

		err = rep.Close()
		require.NoError(t, err)

		err = testDB.Close(ctx)
		require.NoError(t, err)
	})
}

func TestDBCreateTable(t *testing.T) {
	t.Run("db create table success", func(t *testing.T) {
		ctx := context.Background()

		testDB, err := testdb.New(ctx)
		require.NoError(t, err)

		cLog, err := logger.Build("debug")
		require.NoError(t, err)

		rep, err := NewDatabaseRep(testDB.URI, cLog)
		require.NoError(t, err)

		err = rep.createTableMetrics(ctx)
		require.NoError(t, err)

		err = rep.Close()
		require.NoError(t, err)

		err = testDB.Close(ctx)
		require.NoError(t, err)
	})
}

func TestDBSave(t *testing.T) {
	t.Run("db save gauge success", func(t *testing.T) {
		ctx := context.Background()

		testDB, err := testdb.New(ctx)
		require.NoError(t, err)

		cLog, err := logger.Build("debug")
		require.NoError(t, err)

		rep, err := NewDatabaseRep(testDB.URI, cLog)
		require.NoError(t, err)

		err = rep.createTableMetrics(ctx)
		require.NoError(t, err)

		var value float64 = 1
		mtrc := model.Metric{
			Value: &value,
			ID:    "GaugeTest",
			MType: "gauge",
		}

		err = rep.Save(ctx, mtrc)
		require.NoError(t, err)

		err = rep.Close()
		require.NoError(t, err)

		err = testDB.Close(ctx)
		require.NoError(t, err)
	})

	t.Run("db update gauge success", func(t *testing.T) {
		ctx := context.Background()

		testDB, err := testdb.New(ctx)
		require.NoError(t, err)

		cLog, err := logger.Build("debug")
		require.NoError(t, err)

		rep, err := NewDatabaseRep(testDB.URI, cLog)
		require.NoError(t, err)

		err = rep.createTableMetrics(ctx)
		require.NoError(t, err)

		var value = 1.0
		mtrc := model.Metric{
			Value: &value,
			ID:    "GaugeTest",
			MType: "gauge",
		}

		err = rep.Save(ctx, mtrc)
		require.NoError(t, err)

		value = 2.0
		mtrc = model.Metric{
			Value: &value,
			ID:    "GaugeTest",
			MType: "gauge",
		}

		err = rep.Save(ctx, mtrc)
		require.NoError(t, err)

		err = rep.Close()
		require.NoError(t, err)

		err = testDB.Close(ctx)
		require.NoError(t, err)
	})
}

func TestDBPing(t *testing.T) {
	t.Run("db ping success", func(t *testing.T) {
		ctx := context.Background()

		testDB, err := testdb.New(ctx)
		require.NoError(t, err)

		cLog, err := logger.Build("debug")
		require.NoError(t, err)

		rep, err := NewDatabaseRep(testDB.URI, cLog)
		require.NoError(t, err)

		require.NoError(t, rep.Ping(ctx))

		err = rep.Close()
		require.NoError(t, err)

		err = testDB.Close(ctx)
		require.NoError(t, err)
	})
}

func TestDBGet(t *testing.T) {
	t.Run("db get success", func(t *testing.T) {
		ctx := context.Background()

		testDB, err := testdb.New(ctx)
		require.NoError(t, err)

		cLog, err := logger.Build("debug")
		require.NoError(t, err)

		rep, err := NewDatabaseRep(testDB.URI, cLog)
		require.NoError(t, err)

		err = rep.createTableMetrics(ctx)
		require.NoError(t, err)

		var value float64 = 1
		gauge := model.Metric{
			Value: &value,
			ID:    "GaugeTest",
			MType: "gauge",
		}

		var delta int64 = 1
		counter := model.Metric{
			Delta: &delta,
			ID:    "CounterTest",
			MType: "counter",
		}

		err = rep.Save(ctx, gauge)
		require.NoError(t, err)

		err = rep.Save(ctx, counter)
		require.NoError(t, err)

		saved := rep.Get(ctx)
		g, ok := saved[gauge.ID]
		require.Equal(t, ok, true)
		require.Equal(t, g, gauge.ValueString())

		c, ok := saved[counter.ID]
		require.Equal(t, ok, true)
		require.Equal(t, c, counter.DeltaString())

		err = rep.Close()
		require.NoError(t, err)

		err = testDB.Close(ctx)
		require.NoError(t, err)
	})
}

func TestDBMassSave(t *testing.T) {
	t.Run("db mass save success", func(t *testing.T) {
		ctx := context.Background()

		testDB, err := testdb.New(ctx)
		require.NoError(t, err)

		cLog, err := logger.Build("debug")
		require.NoError(t, err)

		rep, err := NewDatabaseRep(testDB.URI, cLog)
		require.NoError(t, err)

		err = rep.createTableMetrics(ctx)
		require.NoError(t, err)

		var value float64 = 1
		gauge := model.Metric{
			Value: &value,
			ID:    "GaugeTest",
			MType: "gauge",
		}

		var delta int64 = 1
		counter := model.Metric{
			Delta: &delta,
			ID:    "CounterTest",
			MType: "counter",
		}

		mtrs := []model.Metric{
			gauge,
			counter,
		}

		err = rep.MassSave(ctx, mtrs)
		require.NoError(t, err)

		saved := rep.Get(ctx)
		g, ok := saved[gauge.ID]
		require.Equal(t, ok, true)
		require.Equal(t, g, gauge.ValueString())

		c, ok := saved[counter.ID]
		require.Equal(t, ok, true)
		require.Equal(t, c, counter.DeltaString())

		err = rep.Close()
		require.NoError(t, err)

		err = testDB.Close(ctx)
		require.NoError(t, err)
	})
}

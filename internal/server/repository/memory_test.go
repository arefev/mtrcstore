package repository

import (
	"context"
	"testing"

	"github.com/arefev/mtrcstore/internal/server/model"
	"github.com/stretchr/testify/require"
)

func TestMemorySave(t *testing.T) {
	t.Run("memory save gauge success", func(t *testing.T) {
		ctx := context.Background()

		var value float64 = 1
		mtrc := model.Metric{
			Value: &value,
			ID:    "GaugeTest",
			MType: "gauge",
		}

		rep := NewMemory()
		err := rep.Save(ctx, mtrc)
		require.NoError(t, err)

		saved, err := rep.findGauge("GaugeTest")
		require.NoError(t, err)
		require.Equal(t, mtrc.Delta, saved.Delta)
	})

	t.Run("memory save counter success", func(t *testing.T) {
		ctx := context.Background()

		var delta int64 = 1
		mtrc := model.Metric{
			Delta: &delta,
			ID:    "CounterTest",
			MType: "counter",
		}

		rep := NewMemory()
		err := rep.Save(ctx, mtrc)
		require.NoError(t, err)

		saved, err := rep.findCounter("CounterTest")
		require.NoError(t, err)
		require.Equal(t, mtrc.Delta, saved.Delta)
	})

	t.Run("memory save gauge error", func(t *testing.T) {
		ctx := context.Background()

		mtrc := model.Metric{
			ID:    "GaugeTest",
			MType: "gauge",
		}

		rep := NewMemory()
		err := rep.Save(ctx, mtrc)
		require.Error(t, err)

		_, err = rep.findGauge("GaugeTest")
		require.Error(t, err)
	})

	t.Run("memory save counter error", func(t *testing.T) {
		ctx := context.Background()

		mtrc := model.Metric{
			ID:    "CounterTest",
			MType: "counter",
		}

		rep := NewMemory()
		err := rep.Save(ctx, mtrc)
		require.Error(t, err)

		_, err = rep.findGauge("CounterTest")
		require.Error(t, err)
	})
}

func TestMemoryFindError(t *testing.T) {
	t.Run("memory find gauge error", func(t *testing.T) {
		ctx := context.Background()

		var delta int64 = 1
		mtrc := model.Metric{
			Delta: &delta,
			ID:    "CounterTest",
			MType: "counter",
		}

		rep := NewMemory()
		err := rep.Save(ctx, mtrc)
		require.NoError(t, err)

		_, err = rep.findGauge("CounterTest")
		require.Error(t, err)
	})

	t.Run("memory find counter error", func(t *testing.T) {
		ctx := context.Background()

		var value float64 = 1
		mtrc := model.Metric{
			Value: &value,
			ID:    "GaugeTest",
			MType: "gauge",
		}

		rep := NewMemory()
		err := rep.Save(ctx, mtrc)
		require.NoError(t, err)

		_, err = rep.findCounter("GaugeTest")
		require.Error(t, err)
	})
}

func TestMemoryFindSuccess(t *testing.T) {
	var delta int64 = 1
	var value = 1.0
	tests := []struct {
		name   string
		metric model.Metric
	}{
		{
			name: "memory find counter success",
			metric: model.Metric{
				Delta: &delta,
				ID:    "CounterTest",
				MType: "counter",
			},
		},
		{
			name: "memory find gauge success",
			metric: model.Metric{
				Value: &value,
				ID:    "GaugeTest",
				MType: "gauge",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			rep := NewMemory()
			err := rep.Save(ctx, tt.metric)
			require.NoError(t, err)

			saved, err := rep.Find(ctx, tt.metric.ID, tt.metric.MType)
			require.NoError(t, err)

			if tt.metric.MType == "counter" {
				require.Equal(t, tt.metric.Delta, saved.Delta)
			} else {
				require.Equal(t, tt.metric.Value, saved.Value)
			}
		})
	}
}

func TestMemoryGet(t *testing.T) {
	t.Run("memory get success", func(t *testing.T) {
		ctx := context.Background()

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

		rep := NewMemory()
		err := rep.Save(ctx, gauge)
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
	})
}

func TestMemoryClose(t *testing.T) {
	t.Run("memory close success", func(t *testing.T) {
		rep := NewMemory()
		require.NoError(t, rep.Close())
	})
}

func TestMemoryPing(t *testing.T) {
	t.Run("memory ping success", func(t *testing.T) {
		ctx := context.Background()
		rep := NewMemory()
		require.NoError(t, rep.Ping(ctx))
	})
}

func TestMemoryMassSave(t *testing.T) {
	t.Run("memory mass save success", func(t *testing.T) {
		ctx := context.Background()

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

		rep := NewMemory()
		err := rep.MassSave(ctx, mtrs)
		require.NoError(t, err)

		saved := rep.Get(ctx)
		g, ok := saved[gauge.ID]
		require.Equal(t, ok, true)
		require.Equal(t, g, gauge.ValueString())

		c, ok := saved[counter.ID]
		require.Equal(t, ok, true)
		require.Equal(t, c, counter.DeltaString())
	})

	t.Run("memory mass save error", func(t *testing.T) {
		ctx := context.Background()

		gauge := model.Metric{
			ID:    "GaugeTest",
			MType: "gauge",
		}

		mtrs := []model.Metric{
			gauge,
		}

		rep := NewMemory()
		err := rep.MassSave(ctx, mtrs)
		require.Error(t, err)

		saved := rep.Get(ctx)
		_, ok := saved[gauge.ID]
		require.Equal(t, ok, false)
	})
}

package repository

import (
	"context"
	"io"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/arefev/mtrcstore/internal/server/logger"
	"github.com/arefev/mtrcstore/internal/server/model"
	"github.com/stretchr/testify/require"
)

func TestFileSave(t *testing.T) {
	t.Run("file save success", func(t *testing.T) {
		ctx := context.Background()
		cLog, err := logger.Build("debug")
		require.NoError(t, err)

		var delta int64 = 1
		mtrc := model.Metric{
			Delta: &delta,
			ID:    "PollCounter",
			MType: "counter",
		}

		rep := NewFile(1, "./storage_test.json", false, cLog)
		err = rep.Save(ctx, mtrc)
		require.NoError(t, err)

		saved, err := rep.findCounter("PollCounter")
		require.NoError(t, err)
		require.Equal(t, mtrc.Delta, saved.Delta)
	})
}

func TestFileWrite(t *testing.T) {
	t.Run("file write success", func(t *testing.T) {
		ctx := context.Background()
		cLog, err := logger.Build("debug")
		require.NoError(t, err)

		var delta int64 = 1
		mtrc := model.Metric{
			Delta: &delta,
			ID:    "PollCounter",
			MType: "counter",
		}

		rep := NewFile(1, "./storage_test.json", false, cLog)
		err = rep.Save(ctx, mtrc)
		require.NoError(t, err)

		rep.write()
		rep.load()

		saved, err := rep.findCounter("PollCounter")
		require.NoError(t, err)
		require.Equal(t, mtrc.Delta, saved.Delta)

		err = os.Remove("./storage_test.json")
		require.NoError(t, err)
	})
}

func TestFileWorker(t *testing.T) {
	t.Run("file worker success", func(t *testing.T) {
		ctx := context.Background()
		const filePermission fs.FileMode = 0o644
		cLog, err := logger.Build("debug")
		require.NoError(t, err)

		var delta int64 = 3
		mtrc := model.Metric{
			Delta: &delta,
			ID:    "PollCounter",
			MType: "counter",
		}

		rep := NewFile(1, "./storage_test.json", false, cLog)
		rep.WorkerRun()
		err = rep.Save(ctx, mtrc)
		require.NoError(t, err)

		time.Sleep(time.Second * 3)

		file, err := os.OpenFile("./storage_test.json", os.O_RDONLY|os.O_CREATE, filePermission)
		require.NoError(t, err)

		data, err := io.ReadAll(file)
		require.NoError(t, err)
		require.Contains(t, string(data), `"PollCounter":3`)

		err = os.Remove("./storage_test.json")
		require.NoError(t, err)
	})
}

func TestFileEvent(t *testing.T) {
	t.Run("file event success", func(t *testing.T) {
		ctx := context.Background()
		const filePermission fs.FileMode = 0o644
		cLog, err := logger.Build("debug")
		require.NoError(t, err)

		var delta int64 = 4
		mtrc := model.Metric{
			Delta: &delta,
			ID:    "PollCounter",
			MType: "counter",
		}

		rep := NewFile(0, "./storage_test.json", false, cLog)
		err = rep.Save(ctx, mtrc)
		require.NoError(t, err)

		file, err := os.OpenFile("./storage_test.json", os.O_RDONLY|os.O_CREATE, filePermission)
		require.NoError(t, err)

		data, err := io.ReadAll(file)
		require.NoError(t, err)
		require.Contains(t, string(data), `"PollCounter":4`)

		err = os.Remove("./storage_test.json")
		require.NoError(t, err)
	})
}

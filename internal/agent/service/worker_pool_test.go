package service_test

import (
	"runtime"
	"testing"

	"github.com/arefev/mtrcstore/internal/agent/repository"
	"github.com/arefev/mtrcstore/internal/agent/service"
	mock_service "github.com/arefev/mtrcstore/internal/agent/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestWorkerPoolRunSuccess(t *testing.T) {
	t.Run("get worker pool run success", func(t *testing.T) {
		var memStats runtime.MemStats

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storage := repository.NewMemory()
		client := mock_service.NewMockSender(ctrl)

		report := service.NewReport(&storage, "http://localhost:8080", "", client)

		runtime.ReadMemStats(&memStats)
		err := report.Save(&memStats)
		require.NoError(t, err)

		mtrs := report.GetMetrics()
		require.NotEmpty(t, mtrs)

		report.IncrementCounter()
		counter := report.Storage.GetCounters()
		require.Equal(t, 1, int(counter["PollCount"]))

		pool := service.NewWorkerPool(report, 3)
		pool.Run()
		pool.Send()

		counter = report.Storage.GetCounters()
		require.Equal(t, 0, int(counter["PollCount"]))
	})
}

package service_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/arefev/mtrcstore/internal/agent/repository"
	"github.com/arefev/mtrcstore/internal/agent/service"
	mock_service "github.com/arefev/mtrcstore/internal/agent/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestWorkerPoolRunSuccess(t *testing.T) {
	t.Run("worker pool run success", func(t *testing.T) {
		ctx := context.Background()
		var memStats runtime.MemStats

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storage := repository.NewMemory()
		client := mock_service.NewMockSender(ctrl)
		client.EXPECT().Request(gomock.Any(), gomock.Any()).MinTimes(1)

		report := service.NewReport(&storage, client)

		runtime.ReadMemStats(&memStats)
		err := report.Save(&memStats)
		require.NoError(t, err)

		mtrs := report.GetMetrics()
		require.NotEmpty(t, mtrs)

		report.IncrementCounter()
		counter := report.Storage.GetCounters()
		require.Equal(t, 1, int(counter["PollCount"]))

		pool := service.NewWorkerPool(report, 3)
		pool.Run(ctx)
		pool.Send()
		time.Sleep(time.Second * 2)

		counter = report.Storage.GetCounters()
		require.Equal(t, 0, int(counter["PollCount"]))
	})
}

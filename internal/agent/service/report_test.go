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

func TestGetGaugesSuccess(t *testing.T) {
	t.Run("get gauges success", func(t *testing.T) {
		var memStats runtime.MemStats

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storage := repository.NewMemory()
		client := mock_service.NewMockSender(ctrl)

		report := service.NewReport(&storage, "http://localhost:8080", "", "", client)

		runtime.ReadMemStats(&memStats)
		err := report.Save(&memStats)
		require.NoError(t, err)

		mtrs := report.GetMetrics()
		require.NotEmpty(t, mtrs)
	})
}

func TestGetCountersSuccess(t *testing.T) {
	t.Run("get counters success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storage := repository.NewMemory()
		client := mock_service.NewMockSender(ctrl)

		report := service.NewReport(&storage, "http://localhost:8080", "", "", client)

		report.IncrementCounter()

		mtrs := report.GetMetrics()
		require.NotEmpty(t, mtrs)
	})
}

func TestSaveCPUSuccess(t *testing.T) {
	t.Run("save cpu success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storage := repository.NewMemory()
		client := mock_service.NewMockSender(ctrl)

		report := service.NewReport(&storage, "http://localhost:8080", "", "", client)

		err := report.SaveCPU()
		require.NoError(t, err)

		mtrs := report.GetMetrics()
		require.NotEmpty(t, mtrs)
	})
}

func TestSendSuccess(t *testing.T) {
	t.Run("send success", func(t *testing.T) {
		var memStats runtime.MemStats

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storage := repository.NewMemory()

		client := mock_service.NewMockSender(ctrl)
		client.EXPECT().DoRequest("http://localhost:8080/updates/", gomock.Any(), gomock.Any())

		report := service.NewReport(&storage, "localhost:8080", "test", "", client)

		runtime.ReadMemStats(&memStats)
		err := report.Save(&memStats)
		require.NoError(t, err)

		err = report.SaveCPU()
		require.NoError(t, err)

		report.IncrementCounter()

		mtrs := report.GetMetrics()
		require.NotEmpty(t, mtrs)

		report.Send(mtrs)
	})
}

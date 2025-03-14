package main

import (
	"context"
	"testing"
	"time"

	"github.com/arefev/mtrcstore/internal/agent"
	mock_service "github.com/arefev/mtrcstore/internal/agent/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestConfigSuccess(t *testing.T) {
	t.Run("test config success", func(t *testing.T) {
		address := "localhost"
		args := []string{
			"-a=" + address,
		}
		conf, err := NewConfig(args)
		require.NoError(t, err)
		require.Equal(t, address, conf.Address)
	})
}

func TestRunSuccess(t *testing.T) {
	t.Run("test run success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock_service.NewMockSender(ctrl)
		client.EXPECT().DoRequest(gomock.Any(), "http://localhost:8080/updates/", gomock.Any(), gomock.Any()).MinTimes(1)

		args := []string{
			"-p=1",
			"-r=2",
		}
		require.ErrorIs(t, run(ctx, args, client), agent.ErrWorkerCanceled)
	})
}

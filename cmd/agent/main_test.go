package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/arefev/mtrcstore/internal/agent"
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 4)
		defer cancel()

		os.Args = []string{
			"-p 1",
		}
		require.ErrorIs(t, run(ctx), agent.ErrWorkerCanceled)
	})
}

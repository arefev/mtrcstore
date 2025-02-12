package main

import (
	"testing"

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

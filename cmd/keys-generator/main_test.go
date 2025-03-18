package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("keys generator run success", func(t *testing.T) {
		const (
			prPath string = "./" + privateKeyName
			puPath string = "./" + publicKeyName
		)
		err := run()
		require.NoError(t, err)
		require.FileExists(t, prPath)
		require.FileExists(t, puPath)

		prData, err := os.ReadFile(prPath)
		require.NoError(t, err)
		require.Contains(t, string(prData), "BEGIN RSA PRIVATE KEY")

		puData, err := os.ReadFile(puPath)
		require.NoError(t, err)
		require.Contains(t, string(puData), "BEGIN RSA PUBLIC KEY")

		require.NoError(t, os.Remove(prPath))
		require.NoError(t, os.Remove(puPath))
	})
}

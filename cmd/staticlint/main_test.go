package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOSExitCheckAnalyzer(t *testing.T) {
	t.Run("staticlint OSExitCheckAnalyzer run success", func(t *testing.T) {
		analysistest.Run(t, analysistest.TestData(), OSExitCheckAnalyzer, "./...")
	})
}

func TestCheckers(t *testing.T) {
	t.Run("staticlint checkers success", func(t *testing.T) {
		checkers := checkers()
		require.NotEmpty(t, checkers)
	})
}

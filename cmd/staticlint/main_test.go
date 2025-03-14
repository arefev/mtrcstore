package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestRun(t *testing.T) {
	t.Run("staticlint run success", func(t *testing.T) {
		analysistest.Run(t, analysistest.TestData(), OSExitCheckAnalyzer, "./...")
	})
}

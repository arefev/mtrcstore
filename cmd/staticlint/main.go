package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"honnef.co/go/tools/quickfix"
)

func main() {
	var checks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
        checks = append(checks, v.Analyzer)
	}

    stChecks := map[string]bool{
        "ST1001": true,
        "ST1003": true,
    }
    for _, v := range stylecheck.Analyzers {
        if stChecks[v.Analyzer.Name] {
            checks = append(checks, v.Analyzer)
        }
	}

    qfChecks  := map[string]bool{
        "QF1001": true,
    }
    for _, v := range quickfix.Analyzers {
        if qfChecks[v.Analyzer.Name] {
            checks = append(checks, v.Analyzer)
        }
	}

	checks = append(
		checks,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	)

	multichecker.Main(
		checks...,
	)
}

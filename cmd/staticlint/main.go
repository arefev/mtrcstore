// Staticlint is multichecker with analyzers
//
//	printf - check consistency of Printf format strings and arguments
//	shadow - check for possible unintended shadowing of variables
//	structtag - checks struct field tags are well formed
//	staticcheck - find bugs and performance issues
//	quickfix - implement code refactorings
//	stylecheck - enforce style rules
//	osexitcheck - check using os.Exit function in main
//
// Steps for run checker:
//  1. Build go build -o ./cmd/staticlint/staticlint ./cmd/staticlint/
//  2. Run ./cmd/staticlint/staticlint ./...
package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

var OSExitCheckAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check using os.Exit function in main",
	Run:  RunOSExitCheckAnalyzer,
}

func main() {
	checks := make([]*analysis.Analyzer, 0)
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

	qfChecks := map[string]bool{
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
		OSExitCheckAnalyzer,
	)

	multichecker.Main(
		checks...,
	)
}

func RunOSExitCheckAnalyzer(pass *analysis.Pass) (any, error) {
	const (
		pkgName      = "os"
		runFuncName  = "Exit"
		fileName     = "main"
		declFuncName = "main"
		reportText   = "os.Exit function using"
	)

	for _, file := range pass.Files {
		fName := ""
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				s, ok := x.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				pkg := fmt.Sprintf("%v", s.X)
				ok = file.Name.String() == fileName && fName == declFuncName && pkg == pkgName && s.Sel.Name == runFuncName
				if !ok {
					return true
				}

				_, ok = x.Args[0].(*ast.BasicLit)
				if !ok {
					return true
				}

				pass.Reportf(x.Pos(), reportText)
			case *ast.FuncDecl:
				fName = x.Name.Name
			}
			return true
		})
	}
	return nil, nil
}

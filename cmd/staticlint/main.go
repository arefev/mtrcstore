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
	Doc:  "check for using os.Exit function",
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
				if s, ok := x.Fun.(*ast.SelectorExpr); ok {
					pkg := fmt.Sprintf("%v", s.X)
					if file.Name.String() == fileName && fName == declFuncName && pkg == pkgName && s.Sel.Name == runFuncName {
						if _, ok := x.Args[0].(*ast.BasicLit); ok {
							pass.Reportf(x.Pos(), reportText)
						}
					}
				}
			case *ast.FuncDecl:
				fName = x.Name.Name
			}
			return true
		})
	}
	return nil, nil
}

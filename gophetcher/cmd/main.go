package main

import "github.com/ole-techwood/AGPWorkbook/internal/analyzer"

func main() {
	analyzerHandler := analyzer.NewAnalyzerHandler()

	analyzerHandler.Audit()
}

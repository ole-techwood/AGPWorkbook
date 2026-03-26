package main

import "github.com/ole-techwood/AGPWorkbook/internal/analysis"

func main() {
	analyzerHandler := analysis.NewAnalysisHandler()

	analyzerHandler.Audit()
}

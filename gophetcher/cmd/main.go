package main

import (
	"os"

	"github.com/ole-techwood/AGPWorkbook/internal/analysis"
)

func main() {
	analyzerHandler := analysis.NewAnalysisHandler()

	if err := analyzerHandler.Audit(); err != nil {
		os.Exit(1)
	}
}

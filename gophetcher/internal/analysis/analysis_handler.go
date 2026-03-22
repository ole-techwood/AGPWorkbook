package analysis

import (
	"fmt"
	"os"

	"github.com/ole-techwood/AGPWorkbook/internal/file"
)

type AnalysisHandler struct {
	fileService     file.FileService
	analysisService AnalysisService
}

func NewAnalysisHandler() AnalysisHandler {
	return AnalysisHandler{
		fileService:     file.NewFileService(),
		analysisService: NewAnalysisService(),
	}
}

func (ah *AnalysisHandler) Audit() {
	filePath, err := ah.fileService.GetTargetFilePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	err = ah.fileService.ValidateFile(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	urls, err := ah.fileService.ReadURLsFromFile(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	auditResults := ah.analysisService.AnalyzeURLs(urls)

	ah.analysisService.PrintResults(auditResults)
}

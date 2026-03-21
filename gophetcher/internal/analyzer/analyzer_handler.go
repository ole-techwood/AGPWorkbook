package analyzer

import (
	"fmt"
	"os"

	"github.com/ole-techwood/AGPWorkbook/internal/file"
)

type AnalyzerHandler struct {
	fileService     *file.FileService
	analyzerService *AnalyzerService
}

func NewAnalyzerHandler() *AnalyzerHandler {
	return &AnalyzerHandler{
		fileService:     file.NewFileService(),
		analyzerService: NewAnalyzerService(),
	}
}

func (ah *AnalyzerHandler) RunAudit() {
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

	err = ah.analyzerService.ValidateURLs(urls)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	auditResults := ah.analyzerService.AnalyzeURLs(urls)

	ah.analyzerService.PrintResults(auditResults)
}

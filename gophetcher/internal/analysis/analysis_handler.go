package analysis

import (
	"fmt"
	"os"
	"strings"

	"github.com/ole-techwood/AGPWorkbook/internal/file"
)

type AnalysisHandler struct {
	fileService file.FileService
	services    map[string]AnalysisService
}

func NewAnalysisHandler() AnalysisHandler {
	webService := NewWebAnalysisService()

	return AnalysisHandler{
		fileService: file.NewFileService(),
		services: map[string]AnalysisService{
			"http":  webService,
			"https": webService,
			"file":  NewFileAnalysisService(),
		},
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

	groups := ah.groupUrlsByService(urls)
	for service, serviceURLs := range groups {
		results := service.AnalyzeURLs(serviceURLs)
		service.PrintResults(results)
	}
}

// getURLBasedService returns the AnalysisService for the scheme of rawURL, or nil if unknown.
func (ah *AnalysisHandler) getURLBasedService(rawURL string) (AnalysisService, bool) {
	scheme, _, found := strings.Cut(rawURL, "://")
	if !found {
		return nil, false
	}

	svc, ok := ah.services[strings.ToLower(scheme)]

	return svc, ok
}

// groupUrlsByService partitions urls into per-service slices, preserving order.
func (ah *AnalysisHandler) groupUrlsByService(urls []string) map[AnalysisService][]string {
	groups := make(map[AnalysisService][]string, len(urls))
	for _, url := range urls {
		service, ok := ah.getURLBasedService(url)
		if !ok {
			continue
		}

		groups[service] = append(groups[service], url)
	}

	return groups
}

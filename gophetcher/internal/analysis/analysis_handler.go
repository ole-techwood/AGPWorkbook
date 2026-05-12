package analysis

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ole-techwood/AGPWorkbook/internal/cli"
	"github.com/ole-techwood/AGPWorkbook/internal/file"
	"github.com/ole-techwood/AGPWorkbook/internal/report"
)

type AnalysisHandler struct {
	cliService  cli.CLIService
	fileService file.FileService
	services    map[string]AnalysisService
}

func NewAnalysisHandler() AnalysisHandler {
	webService := NewWebAnalysisService()

	return AnalysisHandler{
		cliService:  cli.NewCLIService(),
		fileService: file.NewFileService(),
		services: map[string]AnalysisService{
			"http":  webService,
			"https": webService,
			"file":  NewFileAnalysisService(),
		},
	}
}

func (ah *AnalysisHandler) Audit() error {
	startTime := time.Now()
	var totalResources int

	// Defer block to print Final Execution Summary
	defer func() {
		ah.printSummary(startTime, totalResources)
	}()

	options, err := ah.cliService.GetCLIOptions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	// Create timeout context from CLI options
	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()

	err = ah.fileService.ValidateFile(options.FilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	urls, err := ah.fileService.ReadURLsFromFile(options.FilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	totalResources = len(urls)

	// Attach a unique request ID to the context for this audit run.
	ctx = SetAuditRunID(ctx, fmt.Sprintf("audit-%d", time.Now().UnixNano()))

	allResults := make([]AuditResult, 0, len(urls))
	groups := ah.groupUrlsByService(urls)
	for service, serviceURLs := range groups {
		results := service.AnalyzeURLs(ctx, serviceURLs)
		allResults = append(allResults, results...)
	}

	reportHandler := report.NewReportHandler(options.Format)
	reportHandler.Report(ToReportRows(allResults))

	return nil
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

// printSummary prints the final execution summary with total resources and elapsed time.
func (ah *AnalysisHandler) printSummary(startTime time.Time, totalResources int) {
	elapsed := time.Since(startTime)

	fmt.Printf("------------------------------------------------------------\n")
	fmt.Printf("[FINAL SUMMARY]\n")
	fmt.Printf("Total Resources: %d\n", totalResources)
	fmt.Printf("Execution Time: %dms\n", elapsed.Milliseconds())
	fmt.Printf("------------------------------------------------------------\n")
}

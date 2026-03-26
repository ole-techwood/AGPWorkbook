package analysis

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ole-techwood/AGPWorkbook/pkg"
)

type FileAnalysisService struct {
	fileURLValidator *pkg.FileURLValidator
}

var _ AnalysisService = (*FileAnalysisService)(nil)

func NewFileAnalysisService() AnalysisService {
	return &FileAnalysisService{
		fileURLValidator: pkg.NewFileURLValidator(),
	}
}

func (as *FileAnalysisService) AnalyzeURLs(urls []string) []AuditResult {
	results := make([]AuditResult, 0, len(urls))

	for _, rawURL := range urls {
		candidate := strings.TrimSpace(rawURL)

		if _, err := as.fileURLValidator.Validate(candidate); err != nil {
			results = append(results, FileAuditResult{
				BaseAuditResult: BaseAuditResult{
					URL:         candidate,
					Status:      "ERROR",
					ErrorReason: err.Error(),
				},
			})
			continue
		}

		parsed, _ := url.Parse(candidate)
		decodedPath, err := url.PathUnescape(parsed.Path)
		if err != nil {
			results = append(results, FileAuditResult{
				BaseAuditResult: BaseAuditResult{
					URL:         candidate,
					Status:      "ERROR",
					ErrorReason: fmt.Sprintf("error: invalid file path escaping: %v", err),
				},
			})
			continue
		}

		localPath := filepath.Clean(decodedPath)
		info, statErr := os.Stat(localPath)
		if statErr != nil {
			status := "ERROR"
			switch {
			case os.IsNotExist(statErr):
				status = "NOT_FOUND"
			case os.IsPermission(statErr):
				status = "FORBIDDEN"
			}

			results = append(results, FileAuditResult{
				BaseAuditResult: BaseAuditResult{
					URL:         candidate,
					Status:      status,
					ErrorReason: statErr.Error(),
				},
			})
			continue
		}

		results = append(results, FileAuditResult{
			BaseAuditResult: BaseAuditResult{
				URL:    candidate,
				Status: "OK",
			},
			Size:        info.Size(),
			Permissions: info.Mode().Perm().String(),
		})
	}

	return results
}

func fileResultSizeColumn(result FileAuditResult) string {
	if result.Status == "OK" {
		return fmt.Sprintf("%d", result.Size)
	}

	return sanitizeErrorMessage(result.ErrorReason)
}

func toFileAuditResult(result AuditResult) (FileAuditResult, bool) {
	if fileResult, ok := result.(FileAuditResult); ok {
		return fileResult, true
	}
	if ptr, ok := result.(*FileAuditResult); ok && ptr != nil {
		return *ptr, true
	}
	return FileAuditResult{}, false
}

func (as *FileAnalysisService) PrintResults(results []AuditResult) {
	if len(results) == 0 {
		fmt.Println("No audit results to display.")
		return
	}

	fmt.Println(strings.Repeat("-", 72))
	fmt.Printf(tableRowFormat, "RESOURCE", "STATUS", "SIZE", "PERMISSIONS")
	fmt.Println(strings.Repeat("-", 72))

	for _, result := range results {
		fileResult, ok := toFileAuditResult(result)
		if !ok {
			fmt.Printf(tableRowFormat, "-", "ERROR", "-", "unexpected result type")
			continue
		}

		size := fileResultSizeColumn(fileResult)
		fmt.Printf(tableRowFormat, fileResult.URL, fileResult.Status, size, fileResult.Permissions)
	}

	fmt.Println(strings.Repeat("-", 72))
	fmt.Printf("Done. %d resources processed.\n", len(results))
}

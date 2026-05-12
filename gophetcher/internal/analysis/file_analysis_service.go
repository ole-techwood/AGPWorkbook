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

		result := RecoverPanic(candidate, func() AuditResult {
			if _, err := as.fileURLValidator.Validate(candidate); err != nil {
				return FileAuditResult{
					BaseAuditResult: BaseAuditResult{
						URL:         candidate,
						Status:      "ERROR",
						ErrorReason: err.Error(),
					},
				}
			}

			parsed, _ := url.Parse(candidate)
			decodedPath, err := url.PathUnescape(parsed.Path)
			if err != nil {
				return FileAuditResult{
					BaseAuditResult: BaseAuditResult{
						URL:         candidate,
						Status:      "ERROR",
						ErrorReason: fmt.Sprintf("error: invalid file path escaping: %v", err),
					},
				}
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

				return FileAuditResult{
					BaseAuditResult: BaseAuditResult{
						URL:         candidate,
						Status:      status,
						ErrorReason: statErr.Error(),
					},
				}
			}

			return FileAuditResult{
				BaseAuditResult: BaseAuditResult{
					URL:    candidate,
					Status: "OK",
				},
				Size:        info.Size(),
				Permissions: info.Mode().Perm().String(),
			}
		})

		results = append(results, result)
	}

	return results
}

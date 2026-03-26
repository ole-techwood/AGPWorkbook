package analysis

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ole-techwood/AGPWorkbook/pkg"
)

type WebAnalysisService struct {
	client          *http.Client
	webURLValidator *pkg.WebURLValidator
}

var _ AnalysisService = (*WebAnalysisService)(nil)

func NewWebAnalysisService() AnalysisService {
	return &WebAnalysisService{
		client:          &http.Client{Timeout: 12 * time.Second},
		webURLValidator: pkg.NewWebURLValidator(),
	}
}

func (as *WebAnalysisService) AnalyzeURLs(urls []string) []AuditResult {
	results := make([]AuditResult, 0, len(urls))

	for _, rawURL := range urls {
		candidate := strings.TrimSpace(rawURL)

		if _, err := as.webURLValidator.Validate(candidate); err != nil {
			results = append(results, WebAuditResult{
				BaseAuditResult: BaseAuditResult{
					URL:         candidate,
					Status:      "ERROR",
					ErrorReason: err.Error(),
				},
			})
			continue
		}

		resp, err := as.client.Get(candidate)
		if err != nil {
			status := "UNREACHABLE"
			errorMsg := err.Error()
			if strings.Contains(errorMsg, "context deadline exceeded") {
				status = "TIMEOUT"
				errorMsg = "Connection timeout"
			}

			results = append(results, WebAuditResult{
				BaseAuditResult: BaseAuditResult{
					URL:         candidate,
					Status:      status,
					ErrorReason: errorMsg,
				},
			})
			continue
		}

		results = append(results, WebAuditResult{
			BaseAuditResult: BaseAuditResult{
				URL:    candidate,
				Status: strconv.Itoa(resp.StatusCode),
			},
			ContentLength: resp.ContentLength,
			Server:        resp.Header.Get("Server"),
		})
		resp.Body.Close()
	}

	return results
}

func toWebAuditResult(result AuditResult) (WebAuditResult, bool) {
	if webResult, ok := result.(WebAuditResult); ok {
		return webResult, true
	}
	if ptr, ok := result.(*WebAuditResult); ok && ptr != nil {
		return *ptr, true
	}
	return WebAuditResult{}, false
}

func formatSizeColumn(result WebAuditResult) string {
	if _, err := strconv.Atoi(result.Status); err == nil {
		if result.ContentLength >= 0 {
			return strconv.FormatInt(result.ContentLength, 10)
		}

		return "-"
	}

	return sanitizeErrorMessage(result.ErrorReason)
}

func (as *WebAnalysisService) PrintResults(results []AuditResult) {
	if len(results) == 0 {
		fmt.Println("No audit results to display.")
		return
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf(tableRowFormat, "RESOURCE", "STATUS", "SIZE", "SERVER")
	fmt.Println(strings.Repeat("-", 60))

	for _, result := range results {
		webResult, ok := toWebAuditResult(result)
		if !ok {
			fmt.Printf(tableRowFormat, "-", "ERROR", "-", "unexpected result type")
			continue
		}

		size := formatSizeColumn(webResult)
		fmt.Printf(tableRowFormat, webResult.URL, webResult.Status, size, webResult.Server)
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("Done. %d resources processed.\n", len(results))
}

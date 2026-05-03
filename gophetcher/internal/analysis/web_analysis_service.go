package analysis

import (
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

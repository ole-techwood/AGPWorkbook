package analysis

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ole-techwood/AGPWorkbook/pkg"
)

type WebAnalysisService struct {
	client          *http.Client
	webURLValidator *pkg.WebURLValidator
	rng             *rand.Rand
}

var _ AnalysisService = (*WebAnalysisService)(nil)

func NewWebAnalysisService() AnalysisService {
	return &WebAnalysisService{
		client:          &http.Client{Timeout: 12 * time.Second},
		webURLValidator: pkg.NewWebURLValidator(),
		rng:             rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (as *WebAnalysisService) AnalyzeURLs(urls []string) []AuditResult {
	results := make([]AuditResult, 0, len(urls))

	for _, rawURL := range urls {
		candidate := strings.TrimSpace(rawURL)

		result := RecoverPanic(candidate, func() AuditResult {
			as.simulateRandomPanic()

			if _, err := as.webURLValidator.Validate(candidate); err != nil {
				return WebAuditResult{
					BaseAuditResult: BaseAuditResult{
						URL:         candidate,
						Status:      "ERROR",
						ErrorReason: err.Error(),
					},
				}
			}

			resp, err := as.client.Get(candidate)
			if err != nil {
				status := "UNREACHABLE"
				errorMsg := err.Error()
				if strings.Contains(errorMsg, "context deadline exceeded") {
					status = "TIMEOUT"
					errorMsg = "Connection timeout"
				}

				return WebAuditResult{
					BaseAuditResult: BaseAuditResult{
						URL:         candidate,
						Status:      status,
						ErrorReason: errorMsg,
					},
				}
			}

			defer resp.Body.Close()

			return WebAuditResult{
				BaseAuditResult: BaseAuditResult{
					URL:    candidate,
					Status: strconv.Itoa(resp.StatusCode),
				},
				ContentLength: resp.ContentLength,
				Server:        resp.Header.Get("Server"),
			}
		})

		results = append(results, result)
	}

	return results
}

func (as *WebAnalysisService) simulateRandomPanic() {
	// Intentional panic for resilience testing.
	if as.rng.Intn(4) == 0 {
		panic("simulated auditor panic")
	}
}

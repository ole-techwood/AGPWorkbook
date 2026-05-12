package analysis

import (
	"context"
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

func (as *WebAnalysisService) AnalyzeURLs(ctx context.Context, urls []string) []AuditResult {
	results := make([]AuditResult, 0, len(urls))

	for _, rawURL := range urls {
		candidate := strings.TrimSpace(rawURL)

		result := RecoverPanic(ctx, candidate, func() AuditResult {
			as.simulateRandomPanic()

			if _, err := as.webURLValidator.Validate(candidate); err != nil {
				return WebAuditResult{
					BaseAuditResult: BaseAuditResult{
						URL:         candidate,
						Status:      "ERROR",
						ErrorReason: err.Error(),
						RequestID:   GetAuditRunID(ctx),
					},
				}
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, candidate, nil)
			if err != nil {
				return WebAuditResult{
					BaseAuditResult: BaseAuditResult{
						URL:         candidate,
						Status:      "ERROR",
						ErrorReason: err.Error(),
						RequestID:   GetAuditRunID(ctx),
					},
				}
			}

			resp, err := as.client.Do(req)
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
						RequestID:   GetAuditRunID(ctx),
					},
				}
			}

			defer resp.Body.Close()

			return WebAuditResult{
				BaseAuditResult: BaseAuditResult{
					URL:       candidate,
					Status:    strconv.Itoa(resp.StatusCode),
					RequestID: GetAuditRunID(ctx),
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

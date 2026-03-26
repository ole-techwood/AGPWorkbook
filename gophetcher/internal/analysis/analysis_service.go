package analysis

import (
	"strings"
)

type AnalysisService interface {
	AnalyzeURLs(urls []string) []AuditResult
	PrintResults(results []AuditResult)
}

const (
	tableRowFormat = "%-24s %-10s %-10s %s\n"
)

func sanitizeErrorMessage(message string) string {
	message = strings.ReplaceAll(message, "\n", " ")
	message = strings.ReplaceAll(message, "\r", " ")

	return strings.Join(strings.Fields(message), " ")
}

// AnalysisServiceRegistry maps URI schemes to their AnalysisService.
type AnalysisServiceRegistry struct {
	services map[string]AnalysisService
}

func NewAnalysisServiceRegistry() *AnalysisServiceRegistry {
	return &AnalysisServiceRegistry{
		services: map[string]AnalysisService{
			"http":  NewWebAnalysisService(),
			"https": NewWebAnalysisService(),
			"file":  NewFileAnalysisService(),
		},
	}
}

// Resolve returns the AnalysisService for the scheme of rawURL, or nil if unknown.
func (r *AnalysisServiceRegistry) Resolve(rawURL string) (AnalysisService, bool) {
	scheme, _, found := strings.Cut(rawURL, "://")
	if !found {
		return nil, false
	}
	svc, ok := r.services[strings.ToLower(scheme)]
	return svc, ok
}

// GroupByService partitions urls into per-service slices, preserving order.
func (r *AnalysisServiceRegistry) GroupByService(urls []string) map[AnalysisService][]string {
	groups := make(map[AnalysisService][]string, len(urls))
	for _, u := range urls {
		svc, ok := r.Resolve(u)
		if !ok {
			continue
		}
		groups[svc] = append(groups[svc], u)
	}
	return groups
}

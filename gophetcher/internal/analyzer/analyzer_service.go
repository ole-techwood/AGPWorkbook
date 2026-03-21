package analyzer

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type AnalyzerService struct{}

func NewAnalyzerService() *AnalyzerService {
	return &AnalyzerService{}
}

func (as *AnalyzerService) ValidateURLs(urls []string) error {
	if len(urls) == 0 {
		return fmt.Errorf("error: no URLs provided")
	}

	for i, rawURL := range urls {
		candidate := strings.TrimSpace(rawURL)
		if candidate == "" {
			return fmt.Errorf("error: URL at line %d is empty", i+1)
		}

		parsed, err := url.ParseRequestURI(candidate)
		if err != nil {
			return fmt.Errorf("error: malformed URL at line %d: %s", i+1, candidate)
		}

		if parsed.Scheme == "" || parsed.Host == "" {
			return fmt.Errorf("error: URL must include scheme and host at line %d: %s", i+1, candidate)
		}

		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return fmt.Errorf("error: unsupported URL scheme at line %d: %s", i+1, parsed.Scheme)
		}
	}

	return nil
}

func (as *AnalyzerService) AnalyzeURLs(urls []string) []AuditResult {
	var results []AuditResult

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			// Skip URLs that fail to connect
			continue
		}
		defer resp.Body.Close()

		result := AuditResult{
			URL:           url,
			StatusCode:    resp.StatusCode,
			ContentLength: resp.ContentLength,
			Server:        resp.Header.Get("Server"),
		}
		results = append(results, result)
	}

	return results
}

func (as *AnalyzerService) PrintResults(results []AuditResult) {
	if len(results) == 0 {
		fmt.Println("No audit results to display.")
		return
	}

	fmt.Printf("%-30s %-12s %-15s %-20s\n", "URL", "Status Code", "Content Length", "Server")
	fmt.Println(strings.Repeat("-", 80))

	for _, result := range results {
		fmt.Printf("%-30s %-12d %-15d %-20s\n", result.URL, result.StatusCode, result.ContentLength, result.Server)
	}
}

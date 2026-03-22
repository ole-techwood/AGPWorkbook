package analysis

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type AnalysisService struct{}

func NewAnalysisService() AnalysisService {
	return AnalysisService{}
}

func (as *AnalysisService) ValidateURL(candidate string) error {
	if candidate == "" {
		return fmt.Errorf("error: URL is empty")
	}

	parsed, err := url.ParseRequestURI(candidate)
	if err != nil {
		return fmt.Errorf("error: malformed URL: %s", candidate)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("error: URL must include scheme and host: %s", candidate)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("error: unsupported URL scheme: %s", parsed.Scheme)
	}

	return nil
}

func (as *AnalysisService) AnalyzeURLs(urls []string) []AuditResult {
	var results []AuditResult

	// Create HTTP client with 12-second timeout
	client := &http.Client{
		Timeout: 12 * time.Second,
	}

	for _, rawURL := range urls {
		candidate := strings.TrimSpace(rawURL)

		// Validate URL format
		if err := as.ValidateURL(candidate); err != nil {
			results = append(results, AuditResult{
				URL:         candidate,
				Status:      "ERROR",
				ErrorReason: err.Error(),
			})
			continue
		}

		// Perform HTTP request
		resp, err := client.Get(candidate)
		if err != nil {
			// Determine error type
			status := "UNREACHABLE"
			errorMsg := err.Error()

			if strings.Contains(errorMsg, "context deadline exceeded") {
				status = "TIMEOUT"
				errorMsg = "Connection timeout"
			}

			result := AuditResult{
				URL:         candidate,
				Status:      status,
				ErrorReason: errorMsg,
			}
			results = append(results, result)
			continue
		}
		defer resp.Body.Close()

		// Successful response
		result := AuditResult{
			URL:           candidate,
			ContentLength: resp.ContentLength,
			Server:        resp.Header.Get("Server"),
			Status:        strconv.Itoa(resp.StatusCode),
		}
		results = append(results, result)
	}

	return results
}

func (as *AnalysisService) PrintResults(results []AuditResult) {
	if len(results) == 0 {
		fmt.Println("No audit results to display.")
		return
	}

	// Print header
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("%-24s %-10s %-10s %s\n", "URL", "STATUS", "SIZE", "SERVER")
	fmt.Println(strings.Repeat("-", 60))

	// Print each result
	for _, result := range results {
		var sizeStr string

		// Determine SIZE column content
		// Check if Status is a numeric HTTP status code (successful response)
		_, err := strconv.Atoi(result.Status)
		if err == nil {
			// This is a successful HTTP response with a numeric status code
			if result.ContentLength >= 0 {
				sizeStr = strconv.FormatInt(result.ContentLength, 10)
			} else {
				sizeStr = "-"
			}
		} else {
			// This is an error case (ERROR, UNREACHABLE, TIMEOUT, etc.)
			// Clean up error message: replace newlines and multiple spaces
			errorMsg := strings.ReplaceAll(result.ErrorReason, "\n", " ")
			errorMsg = strings.ReplaceAll(errorMsg, "\r", " ")
			errorMsg = strings.Join(strings.Fields(errorMsg), " ")
			sizeStr = errorMsg
		}

		fmt.Printf("%-24s %-10s %-10s %s\n", result.URL, result.Status, sizeStr, result.Server)
	}

	// Print footer
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("Done. %d resources processed.\n", len(results))
}

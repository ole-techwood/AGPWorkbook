# Go Context Examples

## Timeouts and Cancellation

`AnalyzeURLs` accepts a `context.Context` and honours its deadline via `http.NewRequestWithContext`, instead of relying solely on the hard-coded `http.Client` timeout:

```go
func (as *WebAnalysisService) AnalyzeURLs(ctx context.Context, urls []string) []AuditResult {
    results := make([]AuditResult, 0, len(urls))

    for _, rawURL := range urls {
        candidate := strings.TrimSpace(rawURL)

        req, err := http.NewRequestWithContext(ctx, http.MethodGet, candidate, nil)
        if err != nil {
            results = append(results, WebAuditResult{
                BaseAuditResult: BaseAuditResult{
                    URL:         candidate,
                    Status:      "ERROR",
                    ErrorReason: err.Error(),
                },
            })
            continue
        }

        resp, err := as.client.Do(req)
        if err != nil {
            // ctx.Err() distinguishes a deadline/cancellation from a network error.
            if ctx.Err() != nil {
                results = append(results, WebAuditResult{
                    BaseAuditResult: BaseAuditResult{
                        URL:         candidate,
                        Status:      "TIMEOUT",
                        ErrorReason: "context deadline exceeded",
                    },
                })
                break
            }
            // ... handle other errors
        }
        _ = resp
    }

    return results
}
```

The caller wraps the audit with a deadline:

```go
func (ah *AnalysisHandler) Audit() error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Pass ctx down to each service call.
    results := ah.services["https"].AnalyzeURLs(ctx, urls)
    _ = results
    return nil
}
```

## Context Values

Use a typed, unexported key to carry audit-scoped metadata without colliding with other packages:

```go
type auditKey int

const auditRunIDKey auditKey = iota

// WithAuditRunID attaches an audit run ID to the context.
func WithAuditRunID(ctx context.Context, runID string) context.Context {
    return context.WithValue(ctx, auditRunIDKey, runID)
}

// AuditRunID retrieves the audit run ID from the context.
// Returns an empty string if no ID was set.
func AuditRunID(ctx context.Context) string {
    id, _ := ctx.Value(auditRunIDKey).(string)
    return id
}
```

Usage in `AnalysisHandler.Audit`:

```go
ctx = WithAuditRunID(ctx, fmt.Sprintf("audit-%d", time.Now().UnixNano()))
// ctx is now passed to services; each service can call AuditRunID(ctx) for logging.
```

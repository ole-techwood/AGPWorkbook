# Go Error Handling Examples

## Defer Statements

### 1. File Resource Cleanup

This pattern is used throughout gophetcher when reading target files. The `defer` statement ensures the file is closed even if an error occurs during reading.

```go
func (fs *FileService) ReadURLsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error: cannot open file: %v", err)
	}
	defer file.Close()

	urls := make([]string, 0, 4)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error: cannot read file: %v", err)
	}

	return urls, nil
}
```

**Key points:**

- Defer is placed immediately after successful resource acquisition
- The file is guaranteed to close, even if scanner.Scan() or scanner.Err() fails
- Cleanup happens at function return, preventing resource leaks

### 2. HTTP Response Body Cleanup

In gophetcher's web analysis service, HTTP response bodies must be closed to release socket connections back to the connection pool.

```go
func (as *WebAnalysisService) AnalyzeURLs(urls []string) []AuditResult {
	results := make([]AuditResult, 0, len(urls))

	for _, rawURL := range urls {
		candidate := strings.TrimSpace(rawURL)

		resp, err := as.client.Get(candidate)
		if err != nil {
			// Handle network errors (timeout, unreachable)
			results = append(results, WebAuditResult{
				BaseAuditResult: BaseAuditResult{
					URL:         candidate,
					Status:      "UNREACHABLE",
					ErrorReason: err.Error(),
				},
			})
			continue
		}

		// Cleanup response body even before processing
		defer resp.Body.Close()

		results = append(results, WebAuditResult{
			BaseAuditResult: BaseAuditResult{
				URL:    candidate,
				Status: strconv.Itoa(resp.StatusCode),
			},
			ContentLength: resp.ContentLength,
			Server:        resp.Header.Get("Server"),
		})
	}

	return results
}
```

**Key points:**

- HTTP response bodies must be closed to prevent connection leaks
- Defer is placed immediately after successful response, before any processing
- Without defer, if processing fails, the connection would leak

### 3. Deferred Cleanup in Handler Orchestration

When managing multiple services and resources in a handler, defer ensures all cleanup happens in reverse order of acquisition.

```go
func (ah *AnalysisHandler) Audit() {
	options, err := ah.cliService.GetCLIOptions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	err = ah.fileService.ValidateFile(options.FilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	urls, err := ah.fileService.ReadURLsFromFile(options.FilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	allResults := make([]AuditResult, 0, len(urls))
	groups := ah.groupUrlsByService(urls)
	for service, serviceURLs := range groups {
		results := service.AnalyzeURLs(serviceURLs)
		allResults = append(allResults, results...)
	}

	reportHandler := report.NewReportHandler(options.Format)
	reportHandler.Report(ToReportRows(allResults))
}
```

**Key points:**

- When multiple resources are acquired in sequence, defer ensures cleanup happens in reverse order
- If an error occurs at any step, previously acquired resources are still cleaned up
- In gophetcher, the file is automatically closed (via defer in ReadURLsFromFile) before exiting

### 4. Deferred Error Wrapping Pattern

Use defer to add contextual information to errors at function boundaries.

```go
func (cs *CLIService) GetCLIOptions() (*CLIOptions, error) {
	filePath := flag.String("file", "", "path to the file with target URLs")
	format := flag.String("format", "text", "output format: text or json")
	flag.Parse()

	if *filePath == "" {
		return nil, fmt.Errorf("error: -file flag is required")
	}

	requestedFormat := strings.ToLower(strings.TrimSpace(*format))
	supportedFormats := map[string]bool{
		"text": true,
		"json": true,
	}

	if !supportedFormats[requestedFormat] {
		return nil, fmt.Errorf("error: unsupported output format: %s", *format)
	}

	return &CLIOptions{
		FilePath: *filePath,
		Format:   requestedFormat,
	}, nil
}
```

**Key points:**

- Errors are checked immediately after operations
- Error context is wrapped with descriptive messages matching gophetcher's convention (`error: ...`)
- Early return prevents downstream operations on invalid state

## Panic and Recover

### 1. Basic Panic Usage (Project Context)

In gophetcher, panic should remain rare. This example shows a strict constructor guard that panics only when a required dependency is unexpectedly missing due to a programming error.

```go
func NewAnalysisHandler() AnalysisHandler {
	webService := NewWebAnalysisService()
	if webService == nil {
		panic("error: failed to initialize web analysis service")
	}

	return AnalysisHandler{
		cliService:  cli.NewCLIService(),
		fileService: file.NewFileService(),
		services: map[string]AnalysisService{
			"http":  webService,
			"https": webService,
			"file":  NewFileAnalysisService(),
		},
	}
}
```

**Key points:**

- This panic protects against an invalid internal state, not user input errors.
- It aligns with the rule to reserve panic for unrecoverable defects.

### 2. Basic Recover Usage in CLI Entry Point

This pattern wraps the audit flow in `cmd/main.go` so unexpected panics are converted into visible CLI failures.

```go
func main() {
	defer func() {
		if recovered := recover(); recovered != nil {
			fmt.Fprintf(os.Stderr, "error: panic recovered in main: %v\n", recovered)
			os.Exit(1)
		}
	}()

	analyzerHandler := analysis.NewAnalysisHandler()
	if err := analyzerHandler.Audit(); err != nil {
		os.Exit(1)
	}
}
```

**Key points:**

- Recover is used inside a deferred function, which is required by Go.
- The panic value is logged with the project-style `error: ...` prefix.

### 3. Handling Panics for Graceful CLI Exit

This example captures panic, logs context, and exits with a non-zero status after attempting to preserve final diagnostics.

```go
func (ah *AnalysisHandler) AuditWithRecovery() (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			fmt.Fprintf(os.Stderr, "error: audit aborted by panic: %v\n", recovered)
			err = fmt.Errorf("error: audit failed due to unrecoverable internal error")
		}
	}()

	return ah.Audit()
}
```

**Key points:**

- A panic is transformed into an error so callers keep a consistent control flow.
- This keeps the CLI behavior predictable while still surfacing a severe internal issue.

### 4. Real-World Application: Panic Isolation in Per-URL Processing

This example isolates a panic to a single target in `WebAnalysisService`, allowing the audit to continue with the remaining URLs.

```go
func (as *WebAnalysisService) analyzeSingleURL(candidate string) (result WebAuditResult) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = WebAuditResult{
				BaseAuditResult: BaseAuditResult{
					URL:         candidate,
					Status:      "ERROR",
					ErrorReason: fmt.Sprintf("error: panic recovered during web analysis: %v", recovered),
				},
			}
		}
	}()

	resp, err := as.client.Get(candidate)
	if err != nil {
		return WebAuditResult{
			BaseAuditResult: BaseAuditResult{
				URL:         candidate,
				Status:      "UNREACHABLE",
				ErrorReason: err.Error(),
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
}
```

**Key points:**

- Recovery is scoped to one URL analysis, so one bad case does not crash the full run.
- The recovered panic is mapped to a normal audit result row, preserving report continuity.

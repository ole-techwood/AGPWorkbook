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

# Gophetcher

## Overview

Gophetcher is a modular CLI utility designed to audit various types of resources.

## Requirements

### **Input Handling**

The application must read a list of URLs from a local `.txt` file (one URL per line).

### **The Auditor Interface**

Define an `AnalysisService` interface with an `AnalyzeURLs()` method. This abstraction allows the tool to process any resource that can be "inspected" and return a result.

### **Resource Polymorphism**

- **WebAnalysisService:** Handles `http://` and `https://` schemes. Performs an HTTP GET request and captures the Status Code, Content-Length, and "Server" header.
- **FileAnalysisService:** Handles `file://` schemes. Checks if the local file exists, retrieves its size, and captures Unix permissions.
- **Validation:** Implement logic to ensure strings are well-formed before attempting audits.

### **Registry/Selection Logic**

Implement logic to determine the correct `AnalysisService` implementation based on the URI scheme (e.g., `http://` vs `file://`).

### **Resource Management**

Use `defer` statements to guarantee the closing of all opened resources.

- Every `http.Response.Body` in the `WebAnalysisService` must be closed using `defer` immediately after the error check.
- Every `os.File` opened by the `FileAnalysisService` or the input parser must be closed using `defer`.

### **Execution Lifecycle Summary**

Implement a `defer` block in the `Audit` method of `AnalysisHandler` function to ensure that a "Final Execution Summary" (displaying the total time elapsed and the number of resources processed) is printed to the console even if the program encounters a logic error during the audit loop.

### **Panic and Recover**

Implement a safety wrapper for the auditing process.

- The tool must use a `recover()` call within a deferred function inside the specific auditor's `AnalyzeURLs` method.
- If an unexpected runtime error (e.g., a nil pointer dereference or an explicit `panic`) occurs during the audit of a specific resource, the program must capture it, log it as a "CRITICAL" in the report, and continue auditing the remaining resources.

### **The Reporter Interface**

Define a `ReportHandler` interface to abstract the output logic:

- **TextReportHandler:** Prints the formatted summary table to the console.
- **JSONReportHandler:** Outputs the audit results as a structured JSON array for programmatic use.
- Print a summary of the audit using the reporter specified by the user.

### **Context Control**

- **Timeouts:** Implement a per-resource timeout using `context.WithTimeout`.
- **Cancellation:** If the context is cancelled or times out, the tool must stop pending operations immediately and report a "TIMEOUT" status.
- **Request-Scoped Values:** Use `context.WithValue` to pass a unique "Request ID" (e.g., a UUID or a simple counter) from the main loop into each Auditor. This ID must be retrieved inside the Auditor and included in the logs or the `AuditResult`.

### **Code Quality:**

Use a modular structure (e.g., separate logic for file parsing, HTTP fetching, and reporting).

- Define clear types/structs to represent a `Task` and an `AuditResult`.
- The main loop should interact with the `AnalysisService` and `ReportHandler` interfaces rather than concrete structs.
- Follow Go idioms for exported/unexported identifiers.

## CLI Interface

The tool should be invoked with a path to a source file using a flag.

```bash
# Example usage: Default table output
./gophetcher -file=targets.txt

# Example usage: JSON output
./gophetcher -file=targets.txt -format=json
```

**Input File Format (`targets.txt`):**

```text
https://google.com
https://github.com
invalid-url-here
http://localhost:8080
```

**Example Output (Text, with Lifecycle Summary):**

```text
------------------------------------------------------------
RESOURCE                STATUS    SIZE      INFO
------------------------------------------------------------
https://google.com      200       15240     Server: gws
file:///etc/hosts       EXISTS    412       Perm: -rw-r--r--
------------------------------------------------------------
[FINAL SUMMARY]
Total Resources: 2
Execution Time: 150ms
------------------------------------------------------------
```

**The tool includes a `-timeout` flag to control execution duration:**

```bash
# Set a 5-second timeout for the entire audit process
./gophetcher -file=targets.txt -timeout=5s
```

**Example Output (with Request IDs):**

```text
----------------------------------------------------------------------
ID    RESOURCE                STATUS      SIZE      INFO
----------------------------------------------------------------------
101   https://google.com      200         15240     Server: gws
102   https://slow-site.com   TIMEOUT     -         Context deadline exceeded
103   file:///etc/hosts       EXISTS      412       Perm: -rw-r--r--
----------------------------------------------------------------------
[FINAL SUMMARY]
Total Resources: 3
Execution Time: 150ms
----------------------------------------------------------------------
```

## Success Criteria

1. The project compiles and runs using only the standard library (`flag`, `net/http`, `os`, `bufio`, etc.).
2. The code explicitly defines and implements the `AnalysisHandler` and `ReportHandler` interfaces.
3. The tool successfully audits both web URLs and local files in a single run.
4. The user can toggle between table and JSON output via the `format` flag.
5. Invalid resources or failed connections do not crash the program but are reported as errors by the respective auditor.
6. The logic is easy to follow, with consistent naming conventions and no "God objects".
7. Verification that no file descriptors or network connections remain open after the tool finishes execution.
8. The summary section appears at the end of every run, regardless of whether individual audits succeeded or failed.
9. Intentionally triggering a panic in one of the auditors (e.g., by accessing a nil pointer) does not crash the CLI; the tool catches the panic and proceeds to the next URL.
10. The `WebAnalysisService` correctly aborts HTTP requests if the timeout specified in the `timeout` flag is reached.
11. The "Request ID" is successfully passed through the context and displayed in the final report.
12. All `cancel()` functions are called via `defer` to ensure no context-related resources are leaked.

## Constraints

- **Libraries:** `net/http` for requests
- `flag` for CLI arguments.
- No external dependencies.

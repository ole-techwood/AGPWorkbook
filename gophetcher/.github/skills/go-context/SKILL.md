---
name: go-context
description: "Go context patterns for deadlines, cancellation, and request-scoped values. Use when implementing context propagation in handlers or services, adding timeouts to HTTP or file operations, passing request-scoped metadata like audit IDs across API boundaries, or preventing goroutine leaks in audit pipelines."
---

# Go Context

## Understanding Context

Context in Golang represents the execution environment of a request or a process. It carries deadlines, cancellation signals, and other request-scoped values across API boundaries and goroutines.

Contexts can be used to enforce timeouts and cancel operations that exceed a certain duration. This is useful for preventing goroutines from hanging indefinitely and consuming excessive resources.

Contexts can carry request-scoped values across API boundaries. This is useful for passing request IDs, user authentication tokens, or other metadata between functions and goroutines.

## Best Practices

- **Pass Context Explicitly**: Always pass context explicitly to functions and goroutines rather than relying on global variables.
- **Use Context for Cancellation**: Use context cancellation to terminate goroutines gracefully and prevent resource leaks.
- **Set Deadlines Appropriately**: Set timeouts and deadlines on contexts to enforce service-level agreements and prevent long-running operations.
- **Avoid Storing Business Logic in Context**: Context should only be used for request-scoped values and cancellation signals, not for business logic or domain-specific data.

## When to Use This Skill

- Adding timeout control to `AnalyzeURLs` in a service (e.g., replacing the hard-coded `http.Client` timeout with a context deadline)
- Propagating cancellation from `Audit()` down through `AnalysisService` implementations
- Passing audit-scoped metadata (e.g., a run ID) between `AnalysisHandler` and per-scheme services

## References

- [Code examples](references/examples.md)

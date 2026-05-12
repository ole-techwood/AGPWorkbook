---
name: go-error-handling
description: "Go error handling patterns and best practices including defer statements, panic recovery, and resource cleanup. Use when implementing robust error handling, managing resource lifecycle, or designing resilient audit and analysis functions."
---

# Go Error Handling

Comprehensive guidance on error handling patterns in Go, with focus on practical patterns used in gophetcher.

## Defer Statements

Defer statements allow developers to schedule a function call to be executed when the surrounding function returns. This is particularly useful for cleanup tasks and resource management.

## Panic and Recover

Panic is a built-in function in Go that interrupts normal control flow in the current goroutine. In this project, panic should only be used for truly unrecoverable situations where continuing execution would produce invalid behavior.

Recover is a built-in function that allows a deferred function to regain control after a panic. It only works inside deferred functions and is most useful at process boundaries (for example, CLI entry points) to fail gracefully.

Best practices for panic and recover in gophetcher:

- Use panic sparingly and prefer returning errors from services and handlers.
- Recover only at well-defined boundaries, such as the CLI entry point.
- Log error details before panicking or recovering so root causes remain visible.

Reference code and usage examples: [examples](./references/examples.md)

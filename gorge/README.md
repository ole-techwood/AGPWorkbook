# Gorge

## Overview

Gorge is a command-line utility for building and testing data pipelines. The program takes a stream of raw data (such as logs, transactions, or game events), processes it through a chain of transformations, and outputs the final result.

This project contains only test files with the expected results. You must implement business logic so all the tests are green. Good luck!

## Requirements

### Generic Pipeline Core

- Implement a generic `Pipeline[In, Out any]` struct that manages the execution flow.
- Provide a clear abstraction for a processing step. You can use the following standard boilerplate as your starting point:

  ```go
  // Stage defines a single transformation step from type In to type Out.
  type Stage[In, Out any] func(In) (Out, error)
  ```

### Type-Safe Stage Chaining

- Enforce strict compile-time type-safety when linking steps by using a functional composition approach.
- Implement a `Link` function that chains exactly two stages together, where the output type of the first matches the input type of the second:

  ```go
  // Link connects two stages together where the output of the first stage (Middle)
  // perfectly matches the expected input of the second stage.
  func LinkIn, Middle, Out any Stage[In, Out] {}
  ```

- Combine multiple stages sequentially `Link(Link(Stage1, Stage2), Stage3)` to form your final pipeline.

### **Built-in Generic Operations:**

- **Map:** A stage function generator that transforms data from type `A` to type `B` using a custom mapper function:

  ```go
  // Map takes a transform function and returns a Stage that converts A to B.
  func MapA, B any B) Stage[A, B] {}
  ```

- **Filter:** A stage function generator that evaluates data of type `A` against a predicate; if the predicate returns false, it should return a specific error (e.g., `var ErrFiltered = errors.New("filtered")`) to signal skipping or dropping the item:

  ```go
  // Filter takes a predicate function and returns a Stage that passes A through
  // or returns an error if the predicate is false.
  func FilterA any bool) Stage[A, A] {}
  ```

### Concrete Pipeline Demonstration

- **Input Data:** Raw unstructured strings representing game events.
- **Stage 1 (Parse):** Map raw string to a strongly-typed `GameEvent` struct:

  ```go
  type GameEvent struct {
      PlayerID int
      Score    int
      Act      string
  }
  ```

  If the string structure is invalid (e.g., missing fields or wrong separator), return a parsing error

- **Stage 2 (Filter):** Evaluate the valid `GameEvent` struct against a generic predicate. Filter out events where `Score <= 0`.
- **Stage 3 (Transform):** Map the remaining valid `GameEvent` structs into a summary display string.

## CLI Interface

`gorge run --stream="<comma_separated_raw_events>"`: Executes the hardcoded game analytics pipeline against the provided input stream:

```bash
# Input details:
# pid:1|score:50 -> Valid and passes filter
# invalid_garbage_log -> Structurally invalid (Stage 1 Parse Error)
# pid:3|score:0 -> Structurally valid, but filtered out because score is 0 (Stage 2 Filtered)
$ gorge run --stream="pid:1|score:50|act:win,invalid_garbage_log,pid:3|score:0|act:cheat"

--- Pipeline Execution Report ---
[Input]  Total Raw Events: 3
[Stage 1] Failed to parse 1 event(s) due to invalid structure.
[Stage 2] Filtered out 1 event(s) due to zero score.
[Stage 3] Final Output:
  - Player #1 achieved score 50
--------------------------------
```

## Success Criteria

- The pipeline must process all items in the stream sequentially, keeping proper type alignment.
- Your codebase must not use `interface{}` or `any` for type casting inside the pipeline stages; everything must resolve safely at compile time via Go type parameters (Generics).
- Invalid syntax in the `--stream` flag should be handled gracefully, reporting parsing errors without panicking.

## Constraints

- **Standard Library Only:** Zero external dependencies (`go.mod` should only contain the Go version directive).
- **No Code Generation:** All generic type safely must be handled purely by Go's native generics system (`[T any]`), not via `go:generate` or text templates.

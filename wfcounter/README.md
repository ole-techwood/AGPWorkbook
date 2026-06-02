# WFCounter

## Overview

Build a command-line tool named `wfcount` that counts word frequencies across multiple `.txt` files using:

- **Worker Pool** (file-level parallelism)
- **Fan-Out/Fan-In** (intra-file parallelism)
- **Pipeline** (text cleaning stages)
- **Cond** (graceful shutdown & final aggregation signaling)
- **Once** (single initialization of shared resource)

The tool reads multiple text files from a directory provided as a command-line argument, distributes the processing across a fixed number of worker goroutines, and aggregates the results into a single output.

## Requirements

### 1. Input

- **Directory Path**: Accept a directory path as a command-line argument (e.g., go run main.go /path/to/directory).
- **File Type**: Process only .txt files in the specified directory (ignore subdirectories and other file types).
- **File Content**: Assume each .txt file contains plain text (UTF-8 encoded). Words are sequences of characters separated by whitespace.

### 2. Worker Pool

- **Fixed Workers**: Use a Worker Pool with a configurable number of workers (default: 4 workers).
- **Task Distribution**: Each worker reads and processes one file at a time, counting the frequency of each word (case-insensitive, e.g., "The" and "the" are the same word).
- **Concurrency Safety**: Ensure thread-safe aggregation of word counts across all files into a single result map.
- **Channel Usage**: Use channels to distribute file paths to workers and collect results.

### 3. **Fan-Out, Fan-In**

- **Per-File Parallel Processing**
  - For each `.txt` file processed by a worker, implement a **Fan-Out, Fan-In** pattern to parallelize word counting within the file.
  - **Fan-Out**: Split the file’s content into chunks (e.g., by lines or fixed-size segments) and distribute these chunks to multiple goroutines (e.g., 2-4 goroutines per file) to count words in parallel.
  - **Fan-In**: Collect word count results from these goroutines into a single word count map for the file, which is then sent to the Worker Pool’s result channel.
- **Purpose**: This enhances performance for larger files by parallelizing word counting within each file, while the Worker Pool handles parallelism across files.
- **Constraints**:
  - Use a fixed number of goroutines per file (e.g., 2, configurable).
  - Ensure thread-safe merging of word counts from the Fan-In step (e.g., using a `Mutex`).
  - Keep the chunk size reasonable to avoid excessive overhead (e.g., split by lines or ~10 KB chunks).

### 4. Pipeline

- **Text Transformation Pipeline**:
  - Before any word counting begins (including before Fan-Out), the raw file content must pass through a **3-stage pipeline** connected by channels:
    1. **Stage 1 – `ToLower`**: Convert all text to lowercase.
    2. **Stage 2 – `RemovePunctuation`**: Strip all non-letter characters from words (e.g., "hello," → "hello", "world!" → "world").
    3. **Stage 3 – `SplitIntoWords`**: Split the cleaned text into individual words (by whitespace) and emit them one by one.
  - Each stage is a separate goroutine that receives from an input channel and sends to an output channel.
  - The output of Stage 3 (stream of cleaned words) is then fed into the existing Fan-Out, Fan-In word-counting logic.
- **Constraints**:
  - The pipeline must be built dynamically inside each worker for every file.
  - Use unbuffered channels between stages.
  - Stages must close their output channel when done (proper shutdown).
  - The final word stream from Stage 3 must be consumable by the Fan-Out workers.

### 5. Cond

- **Add a clean, graceful shutdown mechanism using `sync.Cond`:**
  - Create a single `sync.Mutex` and a `sync.Cond`

    ```go
    doneCond := sync.NewCond(&mu))
    ```

  - Keep a counter `activeWorkers int` and protect it by the Mutex.
  - Every worker must:
    - Increment `activeWorkers` when it starts processing a file.
    - Decrement `activeWorkers` and call `doneCond.Broadcast()` when it finishes its last task and the task channel is closed.
  - In main:
    - Wait efficiently until all workers are done with:

      ```go
      mu.Lock()
      for activeWorkers > 0 {
          doneCond.Wait()  // releases mu, waits, reacquires mu
      }
      mu.Unlock()
      ```

  - Only after this wait is complete should the program merge the final results and print the sorted output.

- **Purpose**:
  - Demonstrate proper use of `sync.Cond` to signal "all work is finished" without busy-looping or relying solely on closing a results channel (which can be problematic when multiple workers send on it).

### 6. Once

- Add a shared, lazily-initialized global resource that must be set up only once, even when accessed concurrently by all workers:
  - Create a package-level `var stopwordsOnce sync.Once`
  - Create a package-level variable, for example:

    ```go
    var stopwords map[string]struct{} // empty struct for set behavior
    ```

  - Each worker (or any goroutine) may want to use this stopword set to filter out common words like "the", "and", "of", etc.
  - The actual loading of stopwords from a file (e.g. `"stopwords.txt"` in the current directory) must be performed using `stopwordsOnce.Do(func() { ... })`.

- **Rules**:
  - The stopwords file is optional — if it doesn’t exist, the set remains empty (no filtering).
  - The loading logic must be protected by `sync.Once` so that:
    - Even if 100 workers try to access `stopwords` at the same time, the file is read and parsed exactly once.
    - All workers safely read the fully initialized map afterward.
    - After loading, workers should skip any word present in `stopwords` when counting.
- **Purpose**:
  - Practice `sync.Once` for safe, one-time initialization of a shared resource in a highly concurrent environment.

## CLI Interface Specs

- **Command**: go run main.go <directory_path> [-w <num_workers>, -c <num_chunks>]
  - <directory_path>: Required argument for the directory containing .txt files.
  - w <num_workers>: Optional flag to specify the number of workers (default: 4).
  - c <num_chunks>: Optional flag to specify the number of goroutines counting the words in files (default: 2).
- **Output**:
  - Print a sorted list (alphabetically) of words and their total counts across all files.
  - Format: File name and one word per line, word: count.
  - Example:

    ```text
    file.txt:
    apple: 10
    banana: 5
    orange: 8
    ```

- **Error Handling**:
  - If the directory path is invalid or contains no .txt files, print an error message to stderr and exit with code 1.
  - If a file cannot be read, skip it, log an error to stderr, and continue processing other files.

## Success Criteria

- **Correctness**: The program correctly counts word frequencies (case-insensitive) across all `.txt` files in the directory.
- **Concurrency**:
  - The Worker Pool processes files concurrently with no data races or deadlocks.
  - The Fan-Out, Fan-In pattern correctly parallelizes word counting within each file, producing accurate results
  - Pipeline stages process text sequentially per file, but run concurrently with other files and with the Fan-Out workers
- **Output**: The output matches the specified format, sorted alphabetically, with accurate counts.
  - The program must not print results until every worker has finished processing all files.
- **Performance**: The program processes files in parallel (via Worker Pool) and words within files in parallel (via Fan-Out, Fan-In), utilizing multiple goroutines effectively.
- **Validation**:
  - Test with a directory containing 5-10 small .txt files (e.g., 1-10 KB each) and at least one larger file (e.g., 500 KB) to verify Fan-Out, Fan-In benefits
  - Verify counts by manually checking a subset of files.
  - Ensure the program handles edge cases (e.g., empty files, missing directory, invalid files).
  - Verify that "Hello,", "hello", and "HELLO" are all counted as "hello”
  - When stopwords.txt exists and contains words like the, and, a, those words must have count 0 in final output.
  - Removing or renaming stopwords.txt must make those words appear with their real counts.

## Constraints

- **Dependencies**: Use the Go standard library
- **Word Definition**: A word is any sequence of characters separated by whitespace. Ignore punctuation (e.g., "hello," and "hello" are the same word).
- **Case Sensitivity**: Treat words case-insensitively.
- **File Size**: Assume files are small (up to 1 MB each) to keep memory usage simple.

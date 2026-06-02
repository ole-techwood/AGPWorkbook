package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	// "runtime/pprof"
	"sync"
	"time"

	internal "github.com/DonAlexandro/go_advanced/internal"
	slogjson "github.com/veqryn/slog-json"
)

type Worker struct {
	jobs          <-chan string
	results       chan<- internal.FileWordFrequency
	errChan       chan<- error
	counters      *int
	mu            *sync.Mutex
	doneCond      *sync.Cond
	activeWorkers *int
}

func (w Worker) work() {
	for filePath := range w.jobs {
		// Count word frequencies in the file
		words, err := internal.CountWordFrequency(filePath, w.counters)

		if err != nil {
			w.errChan <- err
			continue
		}

		// Create result using the struct
		result := internal.FileWordFrequency{
			FileName: filepath.Base(filePath),
			Words:    words,
		}

		w.results <- result
	}

	// Signal completion only after all channel operations are done
	w.mu.Lock()
	(*w.activeWorkers)--
	w.doneCond.Broadcast()
	w.mu.Unlock()
}

func main() {
	h := slogjson.NewHandler(os.Stdout, &slogjson.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelInfo,
		ReplaceAttr: nil, // Same signature and behavior as stdlib JSONHandler
	})
	// Default global logger
	slog.SetDefault(slog.New(h))

	// Command line flags
	workers := flag.Int("w", 4, "Number of workers to process files concurrently")
	counters := flag.Int("c", 2, "Number of goroutines counting the words in files")

	flag.Parse()

	// Setup profiling
	currentTime := time.Now()
	// profilesDir := "profiles"

	// Start CPU profiling
	// cleanupCPU, err := internal.StartCPUProfiling(profilesDir, currentTime)
	// if err != nil {
	// 	slog.Error("failed to start CPU profiling", slog.Any("error", err))
	// 	os.Exit(1)
	// }
	// defer cleanupCPU()

	// Setup memory profiling
	// defer internal.SetupMemoryProfiling(profilesDir, currentTime)()

	// Cap worker count at maximum of 10
	if *workers > 10 {
		slog.Warn("worker count capped at maximum", slog.Int("requested", *workers), slog.Int("actual", 10))
		*workers = 10
	}

	// Get positional arguments (non-flag arguments)
	args := flag.Args()

	// Check if directory path is provided
	if len(args) == 0 {
		slog.Error("directory path is required")
		slog.Error(fmt.Sprintf("usage: %s <directory_path> [-w <num_workers>]\n", os.Args[0]))
		slog.Error(fmt.Sprintf("example: %s /path/to/files -w 8\n", os.Args[0]))
		os.Exit(1)
	}

	directoryPath := args[0]

	// Validate worker count
	if *workers < 1 {
		slog.Error(fmt.Sprintf("number of workers must be positive, got: %d", *workers))
		os.Exit(1)
	}

	// Validate counter goroutine count
	if *counters < 1 {
		slog.Error(fmt.Sprintf("number of counters must be positive, got: %d", *counters))
		os.Exit(1)
	}

	// Get all .txt files from the directory
	txtFiles, err := internal.GetTxtFiles(directoryPath)
	if err != nil {
		slog.Error("error reading directory", slog.Any("error", err))
		os.Exit(1)
	}

	jobsNum := len(txtFiles)

	jobs := make(chan string, jobsNum)
	results := make(chan internal.FileWordFrequency, jobsNum)
	errChan := make(chan error, jobsNum)

	// Setup graceful shutdown coordination with sync.Cond
	var mu sync.Mutex
	doneCond := sync.NewCond(&mu)
	var activeWorkers int

	var wg sync.WaitGroup

	// Create and start the worker pool
	for w := 1; w <= *workers; w++ {
		// Increment active worker count before spawning
		mu.Lock()
		activeWorkers++
		mu.Unlock()

		wg.Go(func() {
			worker := Worker{
				jobs:          jobs,
				results:       results,
				errChan:       errChan,
				counters:      counters,
				mu:            &mu,
				doneCond:      doneCond,
				activeWorkers: &activeWorkers,
			}

			worker.work()
		})
	}

	// Allow all workers to be spawned before starting job distribution
	// time.Sleep(10 * time.Millisecond)

	// Capture goroutine profile while workers are spawned and waiting for jobs
	// goroutineProfilePath := filepath.Join(profilesDir, fmt.Sprintf("goroutine_profile_%s.prof", currentTime.Format("2006-01-02_15-04-05")))
	// goroutineProfileFile, err := os.Create(goroutineProfilePath)
	// if err != nil {
	// 	slog.Error("failed to create goroutine profile", slog.Any("error", err))
	// } else {
	// 	if err := pprof.Lookup("goroutine").WriteTo(goroutineProfileFile, 0); err != nil {
	// 		slog.Error("failed to write goroutine profile", slog.Any("error", err))
	// 	}
	// 	goroutineProfileFile.Close()
	// 	slog.Info("goroutine profile captured", slog.String("path", goroutineProfilePath))
	// }

	// Send all file paths to jobs channel
	for _, filePath := range txtFiles {
		jobs <- filePath
	}

	// Close the jobs to indicate no more filePaths will be provided after the loop
	close(jobs)

	// Wait for all workers to be spawned and registered
	wg.Wait()

	// Wait for all workers to complete their tasks using sync.Cond
	mu.Lock()
	for activeWorkers > 0 {
		doneCond.Wait() // Releases mu, waits for Broadcast, reacquires mu
	}
	mu.Unlock()

	// Close results and error channels after all workers finished
	close(results)
	close(errChan)

	// Create results directory if it doesn't exist
	resultsDir := "results"
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		slog.Error("failed to create results directory", slog.Any("error", err))
		os.Exit(1)
	}

	// Generate filename with current date
	filename := filepath.Join(resultsDir, fmt.Sprintf("result_%s.md", currentTime.Format("2006-01-02_15-04-05")))

	// Create the output file
	file, err := os.Create(filename)
	if err != nil {
		slog.Error("failed to create output file", slog.Any("error", err))
		os.Exit(1)
	}
	defer file.Close()

	// Collect and write results to file
	for result := range results {
		if _, err := file.WriteString(result.ToHumanReadable()); err != nil {
			slog.Error("failed to write result to file", slog.Any("error", err))
		}
	}

	// Check for any errors
	for err := range errChan {
		if err != nil {
			slog.Error("error processing file", slog.Any("error", err))
		}
	}

	slog.Info("results written to file", slog.String("filename", filename))
}

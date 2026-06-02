package internal

import (
	"os"
	"strings"
	"sync"

	"github.com/mdobak/go-xerrors"
)

type Frequency = map[string]int

type Word struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// ChunkProcessor represents a chunk of text to be processed
type ChunkProcessor struct {
	chunk string
	id    int
}

// ChunkResult represents the result from processing a chunk
type ChunkResult struct {
	frequency map[string]int
	id        int
}

// CountWordFrequency reads a file and counts the frequency of each word using Fan-Out/Fan-In pattern
func CountWordFrequency(filePath string, counters *int) ([]Word, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, xerrors.Newf("failed to read a file %q: %w", filePath, err)
	}

	text := string(content)

	// If text is too small or we only have 1 counter, process sequentially
	if len(text) < 100 || *counters <= 1 {
		frequency := countWordFrequencyInChunk(text)
		return convertFrequencyToWord(frequency), nil
	}

	// Fan-Out: Split text into chunks by lines for better word boundary handling
	numCounters := *counters

	// Create channels for Fan-Out/Fan-In
	jobs := make(chan ChunkProcessor, numCounters)
	results := make(chan ChunkResult, numCounters)

	// Start worker goroutines
	var wg sync.WaitGroup

	for range numCounters {
		wg.Go(func() {
			for job := range jobs {
				frequency := countWordFrequencyInChunk(job.chunk)

				results <- ChunkResult{
					frequency: frequency,
					id:        job.id,
				}
			}
		})
	}

	// Fan-Out: Distribute chunks to workers
	go func() {
		defer close(jobs)

		for i := range numCounters {
			chunk := getTextChunk(text, i, numCounters)

			if chunk == "" {
				continue
			}

			jobs <- ChunkProcessor{
				chunk: chunk,
				id:    i,
			}
		}
	}()

	// Wait for all workers to complete
	wg.Wait()

	// Close the results to indicate no more chunks will be provided
	close(results)

	// Fan-In: Collect and merge results with thread-safe operation
	finalFrequency := mergeChunkFrequenciesIntoSingleFrequency(results)

	return convertFrequencyToWord(finalFrequency), nil
}

func mergeChunkFrequenciesIntoSingleFrequency(results chan ChunkResult) Frequency {
	frequency := make(Frequency)
	var mu sync.Mutex

	for result := range results {
		mu.Lock()
		// Merge frequency maps
		for word, count := range result.frequency {
			frequency[word] += count
		}
		mu.Unlock()
	}

	return frequency
}

// convertFrequencyToWord converts a frequency map to a slice of Word structs
func convertFrequencyToWord(frequency Frequency) []Word {
	result := make([]Word, 0, len(frequency))
	for word, count := range frequency {
		result = append(result, Word{
			Word:  word,
			Count: count,
		})
	}
	return result
}

// countWordFrequencyInChunk processes a text chunk using pipeline and returns word frequencies
func countWordFrequencyInChunk(chunk string) map[string]int {
	// Create preprocessor - pipeline stages will use the package-level sync.Pool
	// to reuse string.Builder instances for reduced memory allocations
	preprocessor := &TextPreprocessor{}
	words := preprocessor.PreprocessText(chunk)

	// Count word frequencies using a map
	frequency := make(Frequency)
	for word := range words {
		frequency[word]++
	}

	return frequency
}

func getChunkSize(lines []string, numCounters int) int {
	// Calculate chunk size (number of lines per chunk)
	chunkSize := len(lines) / numCounters

	return max(chunkSize, 1)
}

func getTextChunk(text string, i, numCounters int) string {
	lines := strings.Split(text, "\n")

	chunkSize := getChunkSize(lines, numCounters)

	start := i * chunkSize
	end := start + chunkSize

	// Handle the last chunk - include remaining lines
	if i == numCounters-1 {
		end = len(lines)
	}

	// Skip empty chunks
	if start >= len(lines) {
		return ""
	}

	// Join lines back into a chunk
	chunk := strings.Join(lines[start:end], "\n")

	return chunk

}

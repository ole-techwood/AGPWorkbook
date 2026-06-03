package pkg

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

var (
	// stopwordsSet stores stopwords for O(1) lookup
	// Empty struct uses zero memory - optimal for set behavior
	stopwordsSet map[string]struct{}

	// stopwordsOnce ensures stopwords are loaded exactly once
	// Thread-safe even with concurrent access from multiple goroutines
	stopwordsOnce sync.Once
)

// loadStopwords reads stopwords from file and initializes the set
// Called via sync.Once - executes exactly once regardless of concurrent calls
func loadStopwords() {
	stopwordsOnce.Do(func() {
		// Initialize empty set - if file doesn't exist, no filtering occurs
		stopwordsSet = make(map[string]struct{})

		// Hardcoded path to stopwords.txt in project root
		file, err := os.Open("stopwords.txt")
		if err != nil {
			// Silent failure - file is optional
			// stopwordsSet remains empty, no filtering will occur
			return
		}
		defer file.Close()

		// Read file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			word := strings.TrimSpace(scanner.Text())
			if word != "" {
				// Store lowercase for case-insensitive matching
				// Words are already lowercased by pipeline before filtering
				stopwordsSet[strings.ToLower(word)] = struct{}{}
			}
		}
		// Ignore scanner errors - partial loading is acceptable
	})
}

// IsStopword checks if a word should be filtered out
// Thread-safe: sync.Once ensures initialization completes before any reads
// After initialization, map is read-only - safe for concurrent access
func IsStopword(word string) bool {
	// Ensure stopwords are loaded before checking
	// sync.Once guarantees this executes exactly once across all goroutines
	loadStopwords()

	// Check if word exists in set
	// Returns false if stopwordsSet is empty (file didn't exist/load failed)
	_, exists := stopwordsSet[word]
	return exists
}

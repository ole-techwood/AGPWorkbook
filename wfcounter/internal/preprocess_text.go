package internal

import (
	"strings"
	"sync"
	"unicode"

	"github.com/DonAlexandro/go_advanced/pkg"
)

// builderPool is a package-level sync.Pool for reusing string.Builder instances
// Shared across all goroutines to maximize buffer reuse efficiency
var builderPool = sync.Pool{
	New: func() any {
		// Pre-allocate builder with 4KB capacity (typical file chunk size)
		// Reduces allocation pressure during text processing
		b := &strings.Builder{}
		b.Grow(65536)
		return b
	},
}

// TextPreprocessor handles text preprocessing using a pipeline pattern
type TextPreprocessor struct{}

// ToLower creates a pipeline stage that converts text to lowercase
func (tp *TextPreprocessor) ToLower(in <-chan string) <-chan string {
	// Create unbuffered output channel
	out := make(chan string)

	// Start goroutine to process input stream asynchronously
	go func() {
		// defer ensures output channel is properly closed
		defer close(out)

		// Process each string from input channel until it's closed
		for text := range in {
			// Get builder from package-level pool to reduce allocations
			builderIface := builderPool.Get()
			builder := builderIface.(*strings.Builder)
			builder.Reset()

			// Manually process runes instead of strings.ToLower to:
			// 1. Use pre-allocated builder from pool
			// 2. Avoid strings.Map's temporary builder allocations
			// 3. Reduce GC pressure
			for _, r := range text {
				builder.WriteRune(unicode.ToLower(r))
			}

			// Extract string and return builder to pool
			lowered := builder.String()
			builder.Reset()
			builderPool.Put(builder)

			// Send lowercased text to output
			out <- lowered
		}
	}()

	// Return receive-only channel immediately
	return out
}

// RemovePunctuation creates a pipeline stage that removes non-alphanumeric characters
// Uses sync.Pool to reuse string.Builder instances across goroutines
func (tp *TextPreprocessor) RemovePunctuation(in <-chan string) <-chan string {
	// Create unbuffered output channel
	out := make(chan string)

	// Start goroutine to process input stream asynchronously
	go func() {
		// defer ensures output channel is properly closed
		defer close(out)

		// Process each string from input channel until it's closed
		for text := range in {
			// Get builder from package-level pool to reduce allocations
			// Avoids allocating new builder for every punctuation removal
			builderIface := builderPool.Get()
			builder := builderIface.(*strings.Builder)
			// Reset clears content while keeping pre-allocated 4KB capacity
			builder.Reset()

			// Manually process runes instead of strings.Map to:
			// 1. Use pre-allocated builder from pool
			// 2. Avoid creating temporary builders in Map's implementation
			// 3. Reduce GC pressure significantly
			for _, r := range text {
				if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
					builder.WriteRune(r)
				} else {
					// Replace punctuation with space to maintain word boundaries
					builder.WriteRune(' ')
				}
			}

			// Extract string while builder is reused
			cleaned := builder.String()

			// Return builder to pool for reuse by other goroutines
			// Reset again to ensure clean state for next user
			builder.Reset()
			builderPool.Put(builder)

			// Send cleaned text to output
			out <- cleaned
		}
		// When input channel closes, defer close(out) signals next stage
	}()

	// Return receive-only channel immediately (non-blocking)
	return out
}

// SplitIntoWords creates a pipeline stage that splits text into individual words
func (tp *TextPreprocessor) SplitIntoWords(in <-chan string) <-chan string {
	// Create unbuffered output channel
	out := make(chan string)

	// Start goroutine to process input stream asynchronously
	go func() {
		// defer ensures output channel is properly closed
		defer close(out)

		// Process each string from input channel until it's closed
		for text := range in {
			// Split text by whitespace into individual words
			words := strings.FieldsSeq(text)

			// Emit each non-empty word to output channel
			for word := range words {
				if len(word) > 0 { // Filter out empty strings
					out <- word
				}
			}
		}
		// When input channel closes and all words emitted,
		// defer close(out) signals consumer that no more words coming
	}()

	// Return receive-only channel immediately (non-blocking)
	return out
}

// FilterStopwords creates a pipeline stage that filters out stopwords
func (tp *TextPreprocessor) FilterStopwords(in <-chan string) <-chan string {
	// Create unbuffered output channel
	out := make(chan string)

	// Start goroutine to process input stream asynchronously
	go func() {
		// defer ensures output channel is properly closed
		defer close(out)

		// Process each word from input channel until it's closed
		for word := range in {
			// Check if word is a stopword using lazy-initialized set
			// sync.Once ensures stopwords load exactly once across all goroutines
			if !pkg.IsStopword(word) {
				// Emit non-stopword to output channel
				out <- word
			}
			// Stopwords are silently filtered out
		}
		// When input channel closes and all words filtered,
		// defer close(out) signals consumer that no more words coming
	}()

	// Return receive-only channel immediately (non-blocking)
	return out
}

// PreprocessText orchestrates the 4-stage pipeline for text preprocessing
func (tp *TextPreprocessor) PreprocessText(text string) <-chan string {
	// Create unbuffered initial channel to feed raw text
	input := make(chan string)

	// Start goroutine to feed the entire text chunk into the pipeline
	go func() {
		// defer ensures channel is closed after sending
		defer close(input)

		// Send the entire text chunk as a single value
		input <- text
		// After sending, defer close(input) executes
		// This signals ToLower stage that no more text is coming
	}()

	// PIPELINE CONSTRUCTION: Chain processing stages together
	// Stage 1: Convert to lowercase
	lowercased := tp.ToLower(input)

	// Stage 2: Remove punctuation from lowercased text
	cleaned := tp.RemovePunctuation(lowercased)

	// Stage 3: Split cleaned text into individual words
	words := tp.SplitIntoWords(cleaned)

	// Stage 4: Filter out stopwords using lazy-initialized set
	filtered := tp.FilterStopwords(words)

	// Return final word stream channel
	// Consumer will receive individual cleaned, non-stopword words one by one
	return filtered
}

package internal

import (
	"fmt"
	"sort"
	"strings"
)

// FileWordFrequency represents the result format for word frequency analysis
type FileWordFrequency struct {
	FileName string `json:"file_name"`
	Words    []Word `json:"words"`
}

// ToHumanReadable converts the struct to human-readable format with sorted words
func (f FileWordFrequency) ToHumanReadable() string {
	// Sort by frequency (descending), then by word (ascending) for ties
	sort.Slice(f.Words, func(i, j int) bool {
		if f.Words[i].Count != f.Words[j].Count {
			return f.Words[i].Count > f.Words[j].Count // Higher frequency first
		}
		return f.Words[i].Word < f.Words[j].Word // Alphabetical order for ties
	})

	// Use strings.Builder for efficient string concatenation
	var builder strings.Builder

	// Pre-allocate capacity for better performance
	// Estimate: filename + ":\n" + (word + ": " + count + "\n") * number of words
	estimatedSize := len(f.FileName) + 2 + len(f.Words)*20 // rough estimate
	builder.Grow(estimatedSize)

	// Build the output string
	builder.WriteString(f.FileName)
	builder.WriteString(":\n")

	for _, w := range f.Words {
		builder.WriteString("\t")
		builder.WriteString(w.Word)
		builder.WriteString(": ")
		builder.WriteString(fmt.Sprintf("%d", w.Count))
		builder.WriteString("\n")
	}

	return builder.String()
}

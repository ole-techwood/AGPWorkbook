package cli

import (
	"fmt"
	"strings"
)

type Report struct {
	TotalRawEvents int
	ParseFailed    int
	FilteredOut    int
	Outputs        []string
}

func (r Report) String() string {
	var b strings.Builder

	b.WriteString("--- Pipeline Execution Report ---\n")

	fmt.Fprintf(&b, "[Input]  Total Raw Events: %d\n", r.TotalRawEvents)
	fmt.Fprintf(&b, "[Stage 1] Failed to parse %d event(s) due to invalid structure.\n", r.ParseFailed)
	fmt.Fprintf(&b, "[Stage 2] Filtered out %d event(s) due to zero score.\n", r.FilteredOut)

	b.WriteString("[Stage 3] Final Output:\n")

	if len(r.Outputs) == 0 {
		b.WriteString("  - (none)\n")
	} else {
		for _, out := range r.Outputs {
			b.WriteString("  - ")
			b.WriteString(out)
			b.WriteByte('\n')
		}
	}
	b.WriteString("--------------------------------\n")

	return b.String()
}

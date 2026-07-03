package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecuteStream_MixedInput(t *testing.T) {
	r := ExecuteStream("pid:1|score:50|act:win,invalid_garbage_log,pid:3|score:0|act:cheat")

	if r.TotalRawEvents != 3 {
		t.Fatalf("total mismatch: got %d want 3", r.TotalRawEvents)
	}
	if r.ParseFailed != 1 {
		t.Fatalf("parse mismatch: got %d want 1", r.ParseFailed)
	}
	if r.FilteredOut != 1 {
		t.Fatalf("filtered mismatch: got %d want 1", r.FilteredOut)
	}
	if len(r.Outputs) != 1 {
		t.Fatalf("outputs len mismatch: got %d want 1", len(r.Outputs))
	}
	if !strings.Contains(strings.ToLower(r.Outputs[0]), "act") {
		t.Fatalf("output must include act: %q", r.Outputs[0])
	}
}

func TestExecuteStream_EmptyTokenCountedAsParseError(t *testing.T) {
	r := ExecuteStream("a,,b")

	if r.TotalRawEvents != 3 {
		t.Fatalf("total mismatch: got %d want 3", r.TotalRawEvents)
	}
	if r.ParseFailed != 3 {
		t.Fatalf("parse mismatch: got %d want 3", r.ParseFailed)
	}
}

func TestExecuteStream_PreservesSuccessfulOutputOrder(t *testing.T) {
	r := ExecuteStream("pid:1|score:10|act:win,pid:2|score:20|act:assist")

	if len(r.Outputs) != 2 {
		t.Fatalf("outputs len mismatch: got %d want 2", len(r.Outputs))
	}
	if !strings.Contains(r.Outputs[0], "Player #1") {
		t.Fatalf("first output order mismatch: %q", r.Outputs[0])
	}
	if !strings.Contains(r.Outputs[1], "Player #2") {
		t.Fatalf("second output order mismatch: %q", r.Outputs[1])
	}
}

func TestRun_MissingStreamFlag_ReturnsNonZero(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	exitCode := Run([]string{}, &out, &errOut)
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code for missing --stream")
	}
	if !strings.Contains(errOut.String(), "missing required flag") {
		t.Fatalf("missing expected error text: %q", errOut.String())
	}
}

func TestReportString_HasExpectedSections(t *testing.T) {
	r := Report{
		TotalRawEvents: 2,
		ParseFailed:    1,
		FilteredOut:    0,
		Outputs:        []string{"Player #1 achieved score 50 with act win"},
	}

	s := r.String()
	for _, part := range []string{
		"Pipeline Execution Report",
		"[Input]",
		"[Stage 1]",
		"[Stage 2]",
		"[Stage 3]",
	} {
		if !strings.Contains(s, part) {
			t.Fatalf("missing report section %q in %q", part, s)
		}
	}
}

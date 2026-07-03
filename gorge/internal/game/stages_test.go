package game

import (
	"strings"
	"testing"
)

func TestToSummary_IncludesAct(t *testing.T) {
	e := GameEvent{PlayerID: 1, Score: 50, Act: "win"}
	out := ToSummary(e)

	if !strings.Contains(out, "Player #1") {
		t.Fatalf("missing player text: %q", out)
	}
	if !strings.Contains(out, "50") {
		t.Fatalf("missing score text: %q", out)
	}
	if !strings.Contains(strings.ToLower(out), "act") || !strings.Contains(strings.ToLower(out), "win") {
		t.Fatalf("missing act text: %q", out)
	}
}

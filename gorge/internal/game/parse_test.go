package game

import "testing"

func TestParseEvent_Valid(t *testing.T) {
	e, err := ParseEvent("pid:1|score:50|act:win")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.PlayerID != 1 || e.Score != 50 || e.Act != "win" {
		t.Fatalf("parsed mismatch: %#v", e)
	}
}

func TestParseEvent_EmptyTokenIsParseError(t *testing.T) {
	_, err := ParseEvent("")
	if err == nil {
		t.Fatal("expected parse error for empty token")
	}
}

func TestParseEvent_InvalidStructure(t *testing.T) {
	_, err := ParseEvent("invalid_garbage_log")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParseEvent_InvalidScore(t *testing.T) {
	_, err := ParseEvent("pid:1|score:abc|act:win")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

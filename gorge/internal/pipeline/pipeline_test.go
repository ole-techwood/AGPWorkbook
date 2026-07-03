package pipeline

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

func TestPipelineRun_ComposedStages(t *testing.T) {
	parse := Stage[string, int](func(in string) (int, error) {
		return strconv.Atoi(in)
	})
	positive := Filter(func(v int) bool { return v > 0 })
	toText := Map(func(v int) string { return "n=" + strconv.Itoa(v) })

	pipe := New(Link(Link(parse, positive), toText))

	out, err := pipe.Run("7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "n=7" {
		t.Fatalf("output mismatch: got %q want %q", out, "n=7")
	}
}

func TestPipelineRun_Filtered(t *testing.T) {
	parse := Stage[string, int](func(in string) (int, error) {
		return strconv.Atoi(in)
	})
	positive := Filter(func(v int) bool { return v > 0 })
	toText := Map(func(v int) string { return strconv.Itoa(v) })

	pipe := New(Link(Link(parse, positive), toText))

	_, err := pipe.Run("0")
	if !errors.Is(err, ErrFiltered) {
		t.Fatalf("expected ErrFiltered, got %v", err)
	}
}

func TestMap_Transforms(t *testing.T) {
	m := Map(func(in string) string {
		return strings.ToUpper(in)
	})

	out, err := m("go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "GO" {
		t.Fatalf("output mismatch: got %q want %q", out, "GO")
	}
}

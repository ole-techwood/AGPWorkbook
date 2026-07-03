package game

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseEvent parses a raw record in format pid:<int>|score:<int>|act:<string>.
func ParseEvent(raw string) (GameEvent, error) {
	if raw == "" {
		return GameEvent{}, fmt.Errorf("invalid event: empty token")
	}

	parts := strings.Split(raw, "|")
	if len(parts) != 3 {
		return GameEvent{}, fmt.Errorf("invalid event: expected 3 fields")
	}

	fields := map[string]string{}
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)

		if len(kv) != 2 || kv[0] == "" {
			return GameEvent{}, fmt.Errorf("invalid event: malformed field")
		}

		if _, exists := fields[kv[0]]; exists {
			return GameEvent{}, fmt.Errorf("invalid event: duplicate field %q", kv[0])
		}

		fields[kv[0]] = kv[1]
	}

	pidRaw, ok := fields["pid"]
	if !ok {
		return GameEvent{}, fmt.Errorf("invalid event: missing pid")
	}

	scoreRaw, ok := fields["score"]
	if !ok {
		return GameEvent{}, fmt.Errorf("invalid event: missing score")
	}

	act, ok := fields["act"]
	if !ok || act == "" {
		return GameEvent{}, fmt.Errorf("invalid event: missing act")
	}

	pid, err := strconv.Atoi(pidRaw)
	if err != nil {
		return GameEvent{}, fmt.Errorf("invalid event: pid parse failed")
	}

	score, err := strconv.Atoi(scoreRaw)
	if err != nil {
		return GameEvent{}, fmt.Errorf("invalid event: score parse failed")
	}

	return GameEvent{PlayerID: pid, Score: score, Act: act}, nil
}

package game

import "fmt"

// ToSummary converts a GameEvent into final display output.
func ToSummary(e GameEvent) string {
	return fmt.Sprintf("Player #%d achieved score %d with act %s", e.PlayerID, e.Score, e.Act)
}

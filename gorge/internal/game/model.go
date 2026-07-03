package game

// GameEvent is structured representation of one raw event record.
type GameEvent struct {
	PlayerID int
	Score    int
	Act      string
}

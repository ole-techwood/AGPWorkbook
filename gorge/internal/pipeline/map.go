package pipeline

// Map lifts a pure transform function into a Stage.
func Map[A, B any](fn func(A) B) Stage[A, B] {
	return func(in A) (B, error) {
		return fn(in), nil
	}
}

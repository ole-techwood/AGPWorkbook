package pipeline

// Link composes two stages with compile-time type alignment.
func Link[In, Middle, Out any](first Stage[In, Middle], second Stage[Middle, Out]) Stage[In, Out] {
	return func(in In) (Out, error) {
		mid, err := first(in)
		if err != nil {
			var zero Out
			return zero, err
		}
		return second(mid)
	}
}

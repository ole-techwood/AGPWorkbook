package pipeline

import "errors"

var ErrFiltered = errors.New("filtered")

// Filter passes a value through when predicate is true, otherwise returns ErrFiltered.
func Filter[A any](predicate func(A) bool) Stage[A, A] {
	return func(in A) (A, error) {
		if !predicate(in) {
			var zero A
			return zero, ErrFiltered
		}
		return in, nil
	}
}

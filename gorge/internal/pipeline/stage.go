package pipeline

// Stage defines one pipeline transformation from In to Out.
type Stage[In, Out any] func(In) (Out, error)

// Pipeline executes a single composed stage.
type Pipeline[In, Out any] struct {
	stage Stage[In, Out]
}

// New constructs a pipeline from a composed stage.
func New[In, Out any](stage Stage[In, Out]) Pipeline[In, Out] {
	return Pipeline[In, Out]{stage: stage}
}

// Run executes the pipeline for one input.
func (p Pipeline[In, Out]) Run(in In) (Out, error) {
	return p.stage(in)
}

package pipeline

import "context"

type StageFunc[I any, O any] func(context.Context, I) (O, error)

type Pipeline[I any, O any] interface {
	Stage(workers int, fn StageFunc[I, O]) Pipeline[I, O]
	Run(ctx context.Context, input <-chan I) (<-chan O, <-chan error)
}

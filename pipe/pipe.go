package pipe

import "github.com/kiriyms/conpats/pool"

// Pool defines the interface for a worker pool that the Pipe uses for concurrent processing.
//
// By default, a pool.Pool instance is used, but custom implementations can be provided via WithPool().
type Pool interface {
	Go(func())
	Wait()
}

type Option func(*Pool)

// WithPool allows specifying a custom Pool implementation for the Pipe to use.
func WithPool(p Pool) Option {
	return func(pool *Pool) {
		*pool = p
	}
}

// PipeFromChan creates a pipe that processes items from the input channel using the provided function and a specified number of workers.
func PipeFromChan[I, O any](fn func(I) O, in <-chan I, workers int, opts ...Option) <-chan O {
	var p Pool = pool.New(workers)
	for _, opt := range opts {
		opt(&p)
	}

	out := make(chan O)

	go func() {
		defer close(out)
		defer p.Wait()
		for item := range in {
			p.Go(func() {
				out <- fn(item)
			})
		}
	}()

	return out
}

// PipeFromSlice creates a pipe that processes items from the input slice using the provided function and a specified number of workers.
func PipeFromSlice[I, O any](fn func(I) O, items []I, workers int, opts ...Option) <-chan O {
	var p Pool = pool.New(workers)
	for _, opt := range opts {
		opt(&p)
	}

	in := make(chan I)
	out := make(chan O)

	go func() {
		for _, item := range items {
			in <- item
		}
		close(in)
	}()

	go func() {
		defer close(out)
		defer p.Wait()
		for item := range in {
			p.Go(func() {
				out <- fn(item)
			})
		}
	}()

	return out
}

// Collect gathers all items from the output channel into a slice.
func Collect[O any](out <-chan O) []O {
	var results []O
	for result := range out {
		results = append(results, result)
	}
	return results
}

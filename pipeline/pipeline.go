package pipeline

import (
	"context"
	"sync"
)

// StageFunc processes input items and sends results to the output channel.
// It should return when ctx is cancelled or input channel is closed.
type StageFunc[In, Out any] func(ctx context.Context, in <-chan In, out chan<- Out)

// Pipeline manages a series of processing stages connected by channels.
// It handles goroutine lifecycle, channel creation, and graceful shutdown.
type Pipeline[T any] struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	input  chan T
}

// New creates a pipeline with the given context and buffer size for channels.
// bufferSize controls how many items can be queued between stages.
func New[T any](ctx context.Context, bufferSize int) *Pipeline[T] {
	// Create a cancellable context so we can shut down all stages
	ctx, cancel := context.WithCancel(ctx)

	return &Pipeline[T]{
		ctx:    ctx,
		cancel: cancel,
		input:  make(chan T, bufferSize),
	}
}

// Input returns the channel where you send data into the pipeline.
func (p *Pipeline[T]) Input() chan<- T {
	return p.input
}

// AddStage adds a processing stage to the pipeline.
// workers controls how many goroutines run this stage concurrently (fan-out).
// Returns the output channel from this stage.
func AddStage[In, Out any](
	p *Pipeline[In],
	workers int,
	bufferSize int,
	fn StageFunc[In, Out],
) *Pipeline[Out] {
	// Create output channel for this stage
	output := make(chan Out, bufferSize)

	// Launch worker goroutines for this stage
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			// Each worker processes items from input and sends to output
			fn(p.ctx, p.input, output)
		}()
	}

	// Close output channel when all workers finish
	go func() {
		p.wg.Wait()
		close(output)
	}()

	// Return new pipeline with this stage's output as input to next stage
	return &Pipeline[Out]{
		ctx:    p.ctx,
		cancel: p.cancel,
		input:  output,
		wg:     sync.WaitGroup{}, // New WaitGroup for next stage
	}
}

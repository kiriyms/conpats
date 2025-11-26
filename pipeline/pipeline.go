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

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

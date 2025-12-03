package pipeline

import "github.com/kiriyms/conpats/pool"

type StageFunc[T any] func(T) T

type WorkerPool interface {
	Go(pool.Job) bool
	CloseAndWait()
}

// Pipeline manages a series of processing stages connected by channels.
// It handles goroutine lifecycle, channel creation, and graceful shutdown.
type Pipeline[T any] struct {
	genOut <-chan T
	stages []StageFunc[T]
	chans  []<-chan T
	pool   WorkerPool
}

// NewFromSlice creates a Pipeline from a slice of integers.
func NewFromSlice[T any](data []T) *Pipeline[T] {
	out := make(chan T)

	go func() {
		defer close(out)
		for _, v := range data {
			out <- v
		}
	}()

	return &Pipeline[T]{genOut: out, stages: make([]StageFunc[T], 0), pool: pool.New(4)}
}

// NewFromChannel creates a Pipeline from an existing channel of integers.
func NewFromChannel[T any](data <-chan T) *Pipeline[T] {
	return &Pipeline[T]{genOut: data, stages: make([]StageFunc[T], 0), pool: pool.New(4)}
}

func (p *Pipeline[T]) AddStage(stage StageFunc[T]) {
	p.stages = append(p.stages, stage)
}

func (p *Pipeline[T]) Run() <-chan T {
	in := p.genOut

	if len(p.stages) == 0 {
		return in
	}

	p.chans = make([]<-chan T, len(p.stages)+1)
	p.chans[0] = p.genOut

	for i, stage := range p.stages {
		out := make(chan T)
		p.pool.Go(func() {
			defer close(out)
			for v := range p.chans[i] {
				out <- stage(v)
			}
		})

		p.chans[i+1] = out
	}

	return p.chans[len(p.chans)-1]
}

func (p *Pipeline[T]) CloseAndWait() {
	p.pool.CloseAndWait()
}

func (p *Pipeline[T]) ConfigurePool(pool WorkerPool) {
	p.pool = pool
}

package pipeline

import "github.com/kiriyms/conpats/pool"

type StageFunc func(int) int

type WorkerPool interface {
	Go(pool.Job) bool
}

// Pipeline manages a series of processing stages connected by channels.
// It handles goroutine lifecycle, channel creation, and graceful shutdown.
type Pipeline struct {
	genOut <-chan int
	stages []StageFunc
	chans  []<-chan int
	pool   WorkerPool
}

// NewFromSlice creates a Pipeline from a slice of integers.
func NewFromSlice(data []int) *Pipeline {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, v := range data {
			out <- v
		}
	}()

	return &Pipeline{genOut: out, stages: make([]StageFunc, 0), pool: pool.New(4)}
}

// NewFromChannel creates a Pipeline from an existing channel of integers.
func NewFromChannel(data <-chan int) *Pipeline {
	return &Pipeline{genOut: data, stages: make([]StageFunc, 0), pool: pool.New(4)}
}

func (p *Pipeline) AddStage(stage StageFunc) {
	p.stages = append(p.stages, stage)
}

func (p *Pipeline) Run() <-chan int {
	in := p.genOut

	if len(p.stages) == 0 {
		return in
	}

	p.chans = make([]<-chan int, len(p.stages)+1)
	p.chans[0] = p.genOut

	for i, stage := range p.stages {
		out := make(chan int)

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

package pipe

import "github.com/kiriyms/conpats/pool"

type Pipe[I, O any] struct {
	in   <-chan I
	out  <-chan O
	fn   func(I) O
	pool *pool.Pool
}

func NewFromChannel[I, O any](fn func(I) O, in <-chan I, workers int) *Pipe[I, O] {
	p := pool.New(workers)
	out := make(chan O)

	go func() {
		defer close(out)
		defer p.CloseAndWait()
		for item := range in {
			p.Go(func() {
				out <- fn(item)
			})
		}
	}()

	return &Pipe[I, O]{
		in:   in,
		out:  out,
		fn:   fn,
		pool: p,
	}
}

func NewFromSlice[I, O any](fn func(I) O, items []I, workers int) *Pipe[I, O] {
	p := pool.New(workers)
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
		defer p.CloseAndWait()
		for item := range in {
			p.Go(func() {
				out <- fn(item)
			})
		}
	}()

	return &Pipe[I, O]{
		in:   in,
		out:  out,
		fn:   fn,
		pool: p,
	}
}

func (p *Pipe[I, O]) Out() <-chan O {
	return p.out
}

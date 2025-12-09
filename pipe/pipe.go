package pipe

import "github.com/kiriyms/conpats/pool"

func PipeFromChan[I, O any](fn func(I) O, in <-chan I, workers int) <-chan O {
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

	return out
}

func PipeFromSlice[I, O any](fn func(I) O, items []I, workers int) <-chan O {
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

	return out
}

func Collect[O any](out <-chan O) []O {
	var results []O
	for result := range out {
		results = append(results, result)
	}
	return results
}

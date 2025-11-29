package pipeline

type StageFunc func(int) int

// Pipeline manages a series of processing stages connected by channels.
// It handles goroutine lifecycle, channel creation, and graceful shutdown.
type Pipeline struct {
	genOut <-chan int
}

func NewFromSlice(data []int) *Pipeline {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, v := range data {
			out <- v
		}
	}()

	return &Pipeline{genOut: out}
}

func NewFromChannel(data <-chan int) *Pipeline {
	return &Pipeline{genOut: data}
}

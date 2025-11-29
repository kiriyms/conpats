package pipeline

type StageFunc func(int) int

// Pipeline manages a series of processing stages connected by channels.
// It handles goroutine lifecycle, channel creation, and graceful shutdown.
type Pipeline struct {
	genOut <-chan int
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

	return &Pipeline{genOut: out}
}

// NewFromChannel creates a Pipeline from an existing channel of integers.
func NewFromChannel(data <-chan int) *Pipeline {
	return &Pipeline{genOut: data}
}

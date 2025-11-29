package pipeline

type StageFunc func(int) int

// Pipeline manages a series of processing stages connected by channels.
// It handles goroutine lifecycle, channel creation, and graceful shutdown.
type Pipeline struct {
}

func NewFromSlice(data []int) *Pipeline {
	return &Pipeline{}
}

func NewFromChannel(data <-chan int) *Pipeline {
	return &Pipeline{}
}

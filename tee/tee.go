package tee

// NewTee takes an input channel and returns n output channels that each receive all items from the input channel.
//
// A buffer size can be specified for the output channels; if buf is 0 or negative, unbuffered channels are created.
func NewTee[I any](in <-chan I, n int, buf int) []chan I {
	if n <= 0 {
		n = 1
	}
	if buf <= 0 {
		buf = 0
	}

	outs := make([]chan I, n)
	for i := range n {
		if buf <= 0 {
			outs[i] = make(chan I)
		} else {
			outs[i] = make(chan I, buf)
		}
	}

	go func() {
		defer func() {
			for _, out := range outs {
				close(out)
			}
		}()

		for item := range in {
			for _, out := range outs {
				out <- item
			}
		}
	}()

	return outs
}

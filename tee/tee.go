package tee

func NewTee[I any](in <-chan I, n int, buf int) []chan I {
	if n <= 0 {
		n = 1
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

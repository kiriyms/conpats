package tee

func Tee[I any](in <-chan I, n int) []<-chan I {
	if n <= 0 {
		n = 1
	}

	outs := make([]chan I, n)
	for i := range n {
		out := make(chan I)
		outs[i] = out
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

	result := make([]<-chan I, n)
	for i, out := range outs {
		result[i] = out
	}
	return result
}

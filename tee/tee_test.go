package tee_test

import (
	"testing"

	"github.com/kiriyms/conpats/tee"
)

func TestTee(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		n := 2

		in := make(chan int)
		go func() {
			defer close(in)
			for i := range 100 {
				in <- i
			}
		}()

		outs := tee.Tee(in, n)

		if len(outs) != n {
			t.Fatalf("expected %d output channels, got %d", n, len(outs))
		}

		results1 := make([]int, 0)
		results2 := make([]int, 0)

		for outs[0] != nil || outs[1] != nil {
			select {
			case v, ok := <-outs[0]:
				if !ok {
					outs[0] = nil
					continue
				}
				results1 = append(results1, v)
			case v, ok := <-outs[1]:
				if !ok {
					outs[1] = nil
					continue
				}
				results2 = append(results2, v)
			}
		}

		if len(results1) != 100 {
			t.Errorf("expected 100 items from first output channel, got %d", len(results1))
		}
		if len(results2) != 100 {
			t.Errorf("expected 100 items from second output channel, got %d", len(results2))
		}
	})
}

package tee_test

import (
	"sync"
	"testing"

	"github.com/kiriyms/conpats/tee"
)

var basicCases = []struct {
	name  string
	chans int
	buf   int
	work  int
}{
	{"basic", 2, 0, 100},
	{"zero channels", 0, 0, 100},
	{"negative channels", -5, 0, 100},
	{"basic buffered", 3, 10, 100},
	{"zero channels buffered", 0, 5, 100},
	{"negative channels buffered", -2, 8, 100},
	{"no work", 2, 0, 0},
}

func TestTee(t *testing.T) {
	t.Parallel()

	for _, tc := range basicCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			in := make(chan int)
			go func() {
				defer close(in)
				for i := range tc.work {
					in <- i
				}
			}()

			outs := tee.NewTee(in, tc.chans, tc.buf)

			if tc.chans > 0 && len(outs) != tc.chans {
				t.Fatalf("expected %d output channels, got %d", tc.chans, len(outs))
			}
			if tc.chans <= 0 && len(outs) != 1 {
				t.Fatalf("expected 1 output channel, got %d", len(outs))
			}

			if tc.buf > 0 {
				for i := range outs {
					if cap(outs[i]) != tc.buf {
						t.Errorf("expected output channel %d to have buffer size %d, got %d", i, tc.buf, cap(outs[i]))
					}
				}
			}

			var results [][]int
			for range outs {
				results = append(results, make([]int, 0))
			}

			var wg sync.WaitGroup
			wg.Add(len(outs))

			for i, out := range outs {
				go func(idx int, ch <-chan int) {
					defer wg.Done()

					for item := range ch {
						results[idx] = append(results[idx], item)
					}
				}(i, out)
			}

			wg.Wait()

			for i := range results {
				if len(results[i]) != tc.work {
					t.Errorf("expected %d items from output channel %d, got %d", tc.work, i, len(results[i]))
					// continue
				}
				for j := range tc.work {
					if results[i][j] != j {
						t.Errorf("expected results[%d][%d] to be %d, got %d", i, j, j, results[i][j])
					}
				}
			}
		})
	}

	t.Run("handles empty input channel", func(t *testing.T) {
		t.Parallel()

		n := 2

		in := make(chan int)
		go func() {
			defer close(in)
		}()

		outs := tee.NewTee(in, n, 0)

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

		if len(results1) != 0 {
			t.Fatalf("expected 0 items from first output channel, got %d", len(results1))
		}
		if len(results2) != 0 {
			t.Fatalf("expected 0 items from second output channel, got %d", len(results2))
		}
		if len(results1) != len(results2) {
			t.Fatalf("expected both output channels to have the same length, got %d and %d", len(results1), len(results2))
		}
	})
}

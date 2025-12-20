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

		if len(results1) != 100 {
			t.Fatalf("expected 100 items from first output channel, got %d", len(results1))
		}
		if len(results2) != 100 {
			t.Fatalf("expected 100 items from second output channel, got %d", len(results2))
		}
		if len(results1) != len(results2) {
			t.Fatalf("expected both output channels to have the same length, got %d and %d", len(results1), len(results2))
		}

		for i := range len(results1) {
			if results1[i] != i {
				t.Errorf("expected results1[%d] to be %d, got %d", i, i, results1[i])
			}
			if results2[i] != i {
				t.Errorf("expected results2[%d] to be %d, got %d", i, i, results2[i])
			}
			if results1[i] != results2[i] {
				t.Errorf("expected results1[%d] to equal results2[%d], got %d and %d", i, i, results1[i], results2[i])
			}
		}
	})

	t.Run("returns 1 channel if input is zero", func(t *testing.T) {
		t.Parallel()

		n := 0

		in := make(chan int)
		go func() {
			defer close(in)
			for i := range 100 {
				in <- i
			}
		}()

		outs := tee.NewTee(in, n, 0)

		if len(outs) != 1 {
			t.Fatalf("expected 1 output channel, got %d", len(outs))
		}

		results := make([]int, 0)

		for item := range outs[0] {
			results = append(results, item)
		}

		if len(results) != 100 {
			t.Fatalf("expected 100 items from output channel, got %d", len(results))
		}

		for i := range len(results) {
			if results[i] != i {
				t.Errorf("expected results[%d] to be %d, got %d", i, i, results[i])
			}
		}
	})

	t.Run("returns 1 channel if input is negative", func(t *testing.T) {
		t.Parallel()

		n := -5

		in := make(chan int)
		go func() {
			defer close(in)
			for i := range 100 {
				in <- i
			}
		}()

		outs := tee.NewTee(in, n, 0)

		if len(outs) != 1 {
			t.Fatalf("expected 1 output channel, got %d", len(outs))
		}

		results := make([]int, 0)

		for item := range outs[0] {
			results = append(results, item)
		}

		if len(results) != 100 {
			t.Fatalf("expected 100 items from output channel, got %d", len(results))
		}

		for i := range len(results) {
			if results[i] != i {
				t.Errorf("expected results[%d] to be %d, got %d", i, i, results[i])
			}
		}
	})

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

	t.Run("handles multiple channels", func(t *testing.T) {
		t.Parallel()

		n := 5

		in := make(chan int)
		go func() {
			defer close(in)
			for i := range 100 {
				in <- i
			}
		}()

		outs := tee.NewTee(in, n, 0)

		if len(outs) != n {
			t.Fatalf("expected %d output channels, got %d", n, len(outs))
		}

		results1 := make([]int, 0)
		results2 := make([]int, 0)
		results3 := make([]int, 0)
		results4 := make([]int, 0)
		results5 := make([]int, 0)

		for outs[0] != nil || outs[1] != nil || outs[2] != nil || outs[3] != nil || outs[4] != nil {
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
			case v, ok := <-outs[2]:
				if !ok {
					outs[2] = nil
					continue
				}
				results3 = append(results3, v)
			case v, ok := <-outs[3]:
				if !ok {
					outs[3] = nil
					continue
				}
				results4 = append(results4, v)
			case v, ok := <-outs[4]:
				if !ok {
					outs[4] = nil
					continue
				}
				results5 = append(results5, v)
			}
		}

		if len(results1) != 100 {
			t.Fatalf("expected 100 items from output channel, got %d", len(results1))
		}
		if len(results2) != 100 {
			t.Fatalf("expected 100 items from output channel, got %d", len(results2))
		}
		if len(results3) != 100 {
			t.Fatalf("expected 100 items from output channel, got %d", len(results3))
		}
		if len(results4) != 100 {
			t.Fatalf("expected 100 items from output channel, got %d", len(results4))
		}
		if len(results5) != 100 {
			t.Fatalf("expected 100 items from output channel, got %d", len(results5))
		}

		for i := range len(results1) {
			if results1[i] != i {
				t.Errorf("expected results1[%d] to be %d, got %d", i, i, results1[i])
			}
			if results2[i] != i {
				t.Errorf("expected results2[%d] to be %d, got %d", i, i, results2[i])
			}
			if results3[i] != i {
				t.Errorf("expected results3[%d] to be %d, got %d", i, i, results3[i])
			}
			if results4[i] != i {
				t.Errorf("expected results4[%d] to be %d, got %d", i, i, results4[i])
			}
			if results5[i] != i {
				t.Errorf("expected results5[%d] to be %d, got %d", i, i, results5[i])
			}
		}
	})

	t.Run("handles buffered channels", func(t *testing.T) {
		t.Parallel()

		n := 3
		buf := 10
		in := make(chan int)
		go func() {
			defer close(in)
			for i := range 50 {
				in <- i
			}
		}()
		outs := tee.NewTee(in, n, buf)

		if len(outs) != n {
			t.Fatalf("expected %d output channels, got %d", n, len(outs))
		}
		for i := range outs {
			if cap(outs[i]) != buf {
				t.Errorf("expected output channel %d to have buffer size %d, got %d", i, buf, cap(outs[i]))
			}
		}

		results := make([][]int, n)
		for i := 0; i < n; i++ {
			results[i] = make([]int, 0)
		}
		for outs[0] != nil || outs[1] != nil || outs[2] != nil {
			select {
			case v, ok := <-outs[0]:
				if !ok {
					outs[0] = nil
					continue
				}
				results[0] = append(results[0], v)
			case v, ok := <-outs[1]:
				if !ok {
					outs[1] = nil
					continue
				}
				results[1] = append(results[1], v)
			case v, ok := <-outs[2]:
				if !ok {
					outs[2] = nil
					continue
				}
				results[2] = append(results[2], v)
			}
		}

		for i := 0; i < n; i++ {
			if len(results[i]) != 50 {
				t.Fatalf("expected 50 items from output channel %d, got %d", i, len(results[i]))
			}
			for j := range len(results[i]) {
				if results[i][j] != j {
					t.Errorf("expected results[%d][%d] to be %d, got %d", i, j, j, results[i][j])
				}
			}
		}
	})
}
package pipe_test

import (
	"fmt"
	"math"
	"sort"
	"testing"

	"github.com/kiriyms/conpats/pipe"
	"github.com/sourcegraph/conc/pool"
)

func TestPipe(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		p := pipe.PipeFromSlice(func(x int) int {
			return x * 2
		}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, 3)

		results := pipe.Collect(p)
		sort.Ints(results)

		expected := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30}
		for i, v := range expected {
			if results[i] != v {
				t.Errorf("expected %d, got %d", v, results[i])
			}
		}
	})

	t.Run("correctly processes from chan", func(t *testing.T) {
		t.Parallel()

		in := make(chan int)
		go func() {
			for i := 1; i <= 10; i++ {
				in <- i
			}
			close(in)
		}()

		p := pipe.PipeFromChan(func(x int) int {
			return x + 1
		}, in, 4)

		results := pipe.Collect(p)
		sort.Ints(results)

		expected := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
		for i, v := range expected {
			if results[i] != v {
				t.Errorf("expected %d, got %d", v, results[i])
			}
		}
	})

	t.Run("correctly handles different worker counts", func(t *testing.T) {
		t.Parallel()

		p := pipe.PipeFromSlice(func(x int) int {
			return x * x
		}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, 4)

		p1 := pipe.PipeFromChan(func(x int) int {
			return x + 1
		}, p, 2)

		p2 := pipe.PipeFromChan(func(x int) int {
			return x * 2
		}, p1, 5)

		results := pipe.Collect(p2)
		sort.Ints(results)

		expected := []int{4, 10, 20, 34, 52, 74, 100, 130, 164, 202, 244, 290, 340, 394, 452}
		for i, v := range expected {
			if results[i] != v {
				t.Errorf("expected %d, got %d", v, results[i])
			}
		}
	})

	t.Run("can custom collect from returned channel", func(t *testing.T) {
		t.Parallel()

		p := pipe.PipeFromSlice(func(x int) int {
			return x + 1
		}, []int{1, 2, 3, 4, 5}, 2)

		p1 := pipe.PipeFromChan(func(x int) string {
			return fmt.Sprintf("Number: %d", x)
		}, p, 2)

		var results []string
		for result := range p1 {
			results = append(results, result)
		}

		sort.Strings(results)

		expected := []string{"Number: 2", "Number: 3", "Number: 4", "Number: 5", "Number: 6"}
		for i, v := range expected {
			if results[i] != v {
				t.Errorf("expected %s, got %s", v, results[i])
			}
		}
	})

	t.Run("handles three-type pipeline", func(t *testing.T) {
		t.Parallel()

		p := pipe.PipeFromSlice(func(x int) float64 {
			return math.Sqrt(float64(x))
		}, []int{1, 4, 9, 16, 25, 36, 49, 64, 81}, 3)

		p1 := pipe.PipeFromChan(func(x float64) string {
			return fmt.Sprintf("Sqrt: %.2f", x)
		}, p, 2)

		results := pipe.Collect(p1)
		sort.Strings(results)

		expected := []string{
			"Sqrt: 1.00",
			"Sqrt: 2.00",
			"Sqrt: 3.00",
			"Sqrt: 4.00",
			"Sqrt: 5.00",
			"Sqrt: 6.00",
			"Sqrt: 7.00",
			"Sqrt: 8.00",
			"Sqrt: 9.00",
		}
		for i, v := range expected {
			if results[i] != v {
				t.Errorf("expected %s, got %s", v, results[i])
			}
		}
	})

	t.Run("returns from int to string in order with 1 worker", func(t *testing.T) {
		t.Parallel()

		p := pipe.PipeFromSlice(func(x int) int {
			return x * x
		}, []int{1, 2, 3, 4, 5}, 1)

		p1 := pipe.PipeFromChan(func(x int) string {
			return fmt.Sprintf("Number: %d", x)
		}, p, 1)

		p2 := pipe.PipeFromChan(func(s string) string {
			return s + "!"
		}, p1, 1)

		results := pipe.Collect(p2)

		expected := []string{"Number: 1!", "Number: 4!", "Number: 9!", "Number: 16!", "Number: 25!"}
		for i, v := range expected {
			if results[i] != v {
				t.Errorf("expected %s, got %s", v, results[i])
			}
		}
	})

	t.Run("handles pipeline with another Pool implementation", func(t *testing.T) {
		t.Parallel()

		pool1 := pool.New()
		p := pipe.PipeFromSlice(func(x int) int {
			return x * x
		}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, 4, pipe.WithPool(pool1))

		p1 := pipe.PipeFromChan(func(x int) int {
			return x + 1
		}, p, 2)

		pool2 := pool.New()
		p2 := pipe.PipeFromChan(func(x int) int {
			return x * 2
		}, p1, 5, pipe.WithPool(pool2))
		results := pipe.Collect(p2)
		sort.Ints(results)

		expected := []int{4, 10, 20, 34, 52, 74, 100, 130, 164, 202, 244, 290, 340, 394, 452}
		for i, v := range expected {
			if results[i] != v {
				t.Errorf("expected %d, got %d", v, results[i])
			}
		}
	})
}

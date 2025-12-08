package pipe_test

import (
	"fmt"
	"testing"

	"github.com/kiriyms/conpats/pipe"
)

func TestPipeline(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		p := pipe.NewFromSlice(func(x int) int {
			return x * x
		}, []int{1, 2, 3, 4, 5}, 1)

		p1 := pipe.NewFromChannel(func(x int) string {
			return fmt.Sprintf("Number: %d", x)
		}, p.Out(), 1)

		p2 := pipe.NewFromChannel(func(s string) string {
			return s + "!"
		}, p1.Out(), 1)

		results := []string{}
		for result := range p2.Out() {
			results = append(results, result)
		}

		expected := []string{"Number: 1!", "Number: 4!", "Number: 9!", "Number: 16!", "Number: 25!"}
		for i, v := range expected {
			if results[i] != v {
				t.Errorf("expected %s, got %s", v, results[i])
			}
		}
	})
}

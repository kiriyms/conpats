package pipeline_test

import (
	"testing"

	"github.com/kiriyms/conpats/pipeline"
)

func TestPipeline(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		pipeline := pipeline.NewFromSlice([]int{1,2,3,4,5})
		pipeline.AddStage(func(i int) int {
			return i * 2
		})
		pipeline.AddStage(func(i int) int {
			return i + 1
		})
		out := pipeline.Run()

		expected := []int{3, 5, 7, 9, 11}
		for _, exp := range expected {
			if got := <-out; got != exp {
				t.Errorf("expected %d, got %d", exp, got)
			}
		}
	})
}

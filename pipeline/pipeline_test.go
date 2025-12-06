package pipeline_test

import (
	"fmt"
	"testing"

	"github.com/kiriyms/conpats/pipeline"
)

func TestPipeline(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		p := pipeline.NewFromSlice(func(x int) string {
			return fmt.Sprintf("Number: %d", x)
		}, []int{1, 2, 3, 4, 5}, 2)

		p1 := pipeline.NewFromChannel(func(s string) string {
			return s + "!"
		}, p.Out(), 2)

		for result := range p1.Out() {
			t.Log(result)
		}
	})
}

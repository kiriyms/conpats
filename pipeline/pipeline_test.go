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

		p := pipeline.New[int]()
		p1 := pipeline.AddStage(p, func(i int) string {
			return fmt.Sprint(i * 2)
		}, 1)
		p2 := pipeline.AddStage(p1, func(s string) float64 {
			var f float64
			fmt.Sscan(s, &f)
			return f / 3.0
		}, 2)
		fmt.Println(p2)
	})
}

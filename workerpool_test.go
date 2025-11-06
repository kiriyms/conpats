package conpats

import "testing"

func TestWorkerPool(t *testing.T) {
	jobs := []int{1, 2, 3, 4, 5}
	f := func(n int) int {
		return n * n
	}
	results := WorkerPool(3, jobs, f)
	for r := range results {
		t.Logf("Result: %d", r)
	}
}
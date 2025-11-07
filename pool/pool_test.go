package pool_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kiriyms/conpats/pool"
)

func TestPool_Basic(t *testing.T) {
	p := pool.New(5)

	for i := 0; i < 10; i++ {
		i := i
		p.Add(func() (any, error) {
			time.Sleep(10 * time.Millisecond)
			return fmt.Sprintf("job-%d", i), nil
		})
	}

	results := p.CloseAndWait()

	if len(results) != 10 {
		t.Fatalf("Expected 10 Results, got %d", len(results))
	}
}

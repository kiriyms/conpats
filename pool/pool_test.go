package pool_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kiriyms/conpats/pool"
)

func TestPool_Basic(t *testing.T) {
	p := pool.New(5)

	for i := range 10 {
		p.Add(func() {
			time.Sleep(2 * time.Second)
			fmt.Printf("job-%d\n", i)
		})
	}

	p.Wait()
}

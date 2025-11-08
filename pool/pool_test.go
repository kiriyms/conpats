package pool_test

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/kiriyms/conpats/pool"
)

func TestPool(t *testing.T) {

	fmt.Println(runtime.NumGoroutine())

	p := pool.New(5)

	fmt.Println(runtime.NumGoroutine())

	for i := range 10 {
		p.Add(func() {
			time.Sleep(2 * time.Second)
			fmt.Printf("job-%d\n", i)
		})
	}

	p.Wait()

	for i := range 10 {
		p.Add(func() {
			time.Sleep(2 * time.Second)
			fmt.Printf("job-%d\n", i)
		})
	}

	p.CloseAndWait()
}

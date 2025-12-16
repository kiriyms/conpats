package pool_test

import (
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kiriyms/conpats/pool"
)

func TestErrorPool(t *testing.T) {
	t.Parallel()

	t.Run("returns 10 errors", func(t *testing.T) {
		t.Parallel()

		p := pool.New(7).WithErrors()
		jobCount := 50
		var completed atomic.Int64
		var errored atomic.Int64

		for i := 0; i < jobCount; i++ {
			p.Go(func() error {
				time.Sleep(2 * time.Millisecond)
				completed.Add(1)

				if i%5 == 0 {
					errored.Add(1)
					return fmt.Errorf("err%d", i)
				}
				return nil
			})
		}

		err := p.Wait().Error()
		err = strings.ReplaceAll(err, "\r\n", "\n")
		errs := strings.Split(err, "\n")
		if completed.Load() != int64(jobCount) {
			t.Errorf("Jobs expected: %d, got: %d", jobCount, completed.Load())
		}
		if len(errs) != int(errored.Load()) {
			t.Errorf("Errors count mismatch; count: %d, collected: %d", errored.Load(), len(errs))
		}
	})

	t.Run("returns nil if no errors", func(t *testing.T) {
		t.Parallel()

		p := pool.New(7).WithErrors()
		jobCount := 50
		var completed atomic.Int64

		for i := 0; i < jobCount; i++ {
			p.Go(func() error {
				time.Sleep(2 * time.Millisecond)
				completed.Add(1)
				return nil
			})
		}

		err := p.Wait()
		if completed.Load() != int64(jobCount) {
			t.Errorf("Jobs expected: %d, got: %d", jobCount, completed.Load())
		}
		if err != nil {
			t.Errorf("Expected err to be nil")
		}
	})
}

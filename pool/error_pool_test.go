package pool_test

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kiriyms/conpats/pool"
)

type ErrorIdContainer struct {
	sync.Mutex
	id int64
}

func TestErrorPool(t *testing.T) {
	t.Parallel()

	t.Run("returns 10 errors", func(t *testing.T) {
		t.Parallel()

		p := pool.New(7).WithErrors(false)
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

		errs := p.Wait()

		sort.Slice(errs, func(i, j int) bool {
			var n int
			var m int
			fmt.Sscanf(errs[i].Error(), "err%d", &n)
			fmt.Sscanf(errs[j].Error(), "err%d", &m)
			return n < m
		})

		if completed.Load() != int64(jobCount) {
			t.Errorf("Jobs expected: %d, got: %d", jobCount, completed.Load())
		}
		if len(errs) != int(errored.Load()) {
			t.Errorf("Errors count mismatch; count: %d, collected: %d", errored.Load(), len(errs))
		}

		checkList := []string{"err0", "err5", "err10", "err15", "err20", "err25", "err30", "err35", "err40", "err45"}
		for i, err := range errs {
			if err.Error() != checkList[i] {
				t.Errorf("Expected error %s, got %s", checkList[i], err.Error())
			}
		}
	})

	t.Run("returns nil if no errors", func(t *testing.T) {
		t.Parallel()

		p := pool.New(7).WithErrors(false)
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

	t.Run("returns only first error", func(t *testing.T) {
		t.Parallel()

		p := pool.New(7).WithErrors(true)
		jobCount := 50
		var completed atomic.Int64
		var errored atomic.Int64

		errorId := ErrorIdContainer{id: -1}

		for i := 0; i < jobCount; i++ {
			p.Go(func() error {
				time.Sleep(2 * time.Millisecond)
				completed.Add(1)

				if i%5 == 0 {
					errorId.Lock()
					defer errorId.Unlock()
					errored.Add(1)
					if errorId.id == -1 {
						errorId.id = int64(i)
					}
					return fmt.Errorf("err%d", i)
				}
				return nil
			})
		}

		errs := p.Wait()

		if completed.Load() != int64(jobCount) {
			t.Errorf("Jobs expected: %d, got: %d", jobCount, completed.Load())
		}
		if len(errs) != 1 {
			t.Errorf("Expected only one error, got %d", len(errs))
		}
		expectedErr := fmt.Sprintf("err%d", errorId.id)
		if errs[0].Error() != expectedErr {
			t.Errorf("Expected error '%s', got '%s'", expectedErr, errs[0].Error())
		}
	})
}

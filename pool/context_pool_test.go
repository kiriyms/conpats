package pool_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kiriyms/conpats/pool"
)

func TestContextPool(t *testing.T) {
	t.Parallel()

	t.Run("collects errors and passes context", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p := pool.New(7).WithErrors(false).WithContext(ctx)
		jobCount := 50

		var completed atomic.Int64
		var errored atomic.Int64
		var ctxUsed atomic.Int64

		for i := 0; i < jobCount; i++ {
			i := i
			p.Go(func(c context.Context) error {
				if c != nil && c != context.Background() {
					ctxUsed.Add(1)
				}

				time.Sleep(1 * time.Millisecond)
				completed.Add(1)

				if i%7 == 0 {
					errored.Add(1)
					return fmt.Errorf("err%d", i)
				}
				return nil
			})
		}

		err := p.Wait()

		if ctxUsed.Load() != int64(jobCount) {
			t.Errorf("Expected all jobs to receive context, got %d/%d", ctxUsed.Load(), jobCount)
		}

		if completed.Load() != int64(jobCount) {
			t.Errorf("Jobs expected: %d, got: %d", jobCount, completed.Load())
		}

		if err == nil {
			t.Fatalf("Expected aggregated error, got nil")
		}

		if len(err) != int(errored.Load()) {
			t.Errorf("Errors count mismatch; expected: %d, collected: %d", errored.Load(), len(err))
		}
	})

	t.Run("returns nil if no errors", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p := pool.New(7).WithErrors(false).WithContext(ctx)

		jobCount := 30
		var completed atomic.Int64

		for i := 0; i < jobCount; i++ {
			p.Go(func(c context.Context) error {
				if c == nil {
					t.Fatalf("Expected context to be passed")
				}

				completed.Add(1)
				return nil
			})
		}

		err := p.Wait()

		if completed.Load() != int64(jobCount) {
			t.Errorf("Jobs expected: %d, got: %d", jobCount, completed.Load())
		}

		if err != nil {
			t.Errorf("Expected error to be nil, got: %v", err)
		}
	})

	t.Run("work cancels on parent context cancel", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p := pool.New(7).WithErrors(false).WithContext(ctx)

		jobCount := 4
		var completed atomic.Int64
		var cancelled atomic.Int64

		for i := 0; i < jobCount; i++ {
			p.Go(func(c context.Context) error {
				select {
				case <-c.Done():
					cancelled.Add(1)
					return c.Err()
				case <-time.After(50 * time.Millisecond):
					completed.Add(1)
					return nil
				}
			})
		}

		cancel()

		err := p.Wait()

		if len(err) != jobCount {
			t.Fatalf("Expected %d errors, got: %d", jobCount, len(err))
		}

		for _, e := range err {
			if e != context.Canceled {
				t.Errorf("Expected context.Canceled error, got: %v", e)
			}
		}

		if completed.Load() != 0 {
			t.Errorf("Expected 0 completed jobs, got: %d", completed.Load())
		}
		if cancelled.Load() != 4 {
			t.Errorf("Expected 4 cancelled jobs, got: %d", cancelled.Load())
		}
	})
}

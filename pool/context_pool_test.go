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

		p := pool.New(7).WithErrors().WithContext(ctx)
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

		p := pool.New(7).WithErrors().WithContext(ctx)

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

		p := pool.New(7).WithErrors().WithContext(ctx)

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

	t.Run("handles TryGo correctly", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p := pool.New(1).WithErrors().WithContext(ctx)
		ok := p.TryGo(func(c context.Context) error {
			time.Sleep(2 * time.Millisecond)
			return nil
		})

		if !ok {
			t.Fatalf("Should not error on .TryGo() if the Pool is not closed yet")
		}

		p.Wait()

		ok = p.TryGo(func(c context.Context) error {
			time.Sleep(2 * time.Millisecond)
			return nil
		})

		if ok {
			t.Fatalf("Expected error on .TryGo() because Pool is closed")
		}
	})

	t.Run("cancels on error correctly", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p := pool.New(5).WithErrors().WithContext(ctx, pool.WithCancelOnErr())
		jobCount := 50
		var completed atomic.Int64
		var cancelled atomic.Int64

		for i := 0; i < jobCount; i++ {
			p.Go(func(c context.Context) error {
				select {
				case <-c.Done():
					cancelled.Add(1)
					return c.Err()
				case <-time.After(2 * time.Millisecond):
					if completed.Load() == 5 {
						cancelled.Add(1)
						return fmt.Errorf("intentional error")
					}
					completed.Add(1)
					return nil
				}
			})
		}

		errs := p.Wait()
		if len(errs) != 45 {
			t.Errorf("Expected 45 errors, got: %d", len(errs))
		}
		if completed.Load() != 5 {
			t.Errorf("Expected 5 completed jobs, got: %d", completed.Load())
		}
		if cancelled.Load() != 45 {
			t.Errorf("Expected 45 cancelled jobs, got: %d", cancelled.Load())
		}
	})

	t.Run("cancels correctly on error and returns only first error", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p := pool.New(5).WithErrors(pool.WithOnlyFirstErr()).WithContext(ctx, pool.WithCancelOnErr())
		jobCount := 50
		var completed atomic.Int64
		var cancelled atomic.Int64

		for i := 0; i < jobCount; i++ {
			p.Go(func(c context.Context) error {
				select {
				case <-c.Done():
					cancelled.Add(1)
					return c.Err()
				case <-time.After(2 * time.Millisecond):
					if completed.Load() == 5 {
						cancelled.Add(1)
						return fmt.Errorf("intentional error")
					}
					completed.Add(1)
					return nil
				}
			})
		}

		errs := p.Wait()
		if len(errs) != 1 {
			t.Errorf("Expected only first error, got: %d", len(errs))
		}
		if completed.Load() != 5 {
			t.Errorf("Expected 5 completed jobs, got: %d", completed.Load())
		}
		if cancelled.Load() != 45 {
			t.Errorf("Expected 45 cancelled jobs, got: %d", cancelled.Load())
		}
		if errs[0].Error() != "intentional error" {
			t.Errorf("Expected 'intentional error', got: %v", errs[0])
		}
	})

	t.Run("handles TryGo correctly with cancel on error", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p := pool.New(1).WithErrors(pool.WithOnlyFirstErr()).WithContext(ctx, pool.WithCancelOnErr())
		ok := p.TryGo(func(c context.Context) error {
			time.Sleep(2 * time.Millisecond)
			return fmt.Errorf("intentional error")
		})

		if !ok {
			t.Fatalf("Should not error on .TryGo() if the Pool is not closed yet")
		}

		p.Wait()

		ok = p.TryGo(func(c context.Context) error {
			time.Sleep(2 * time.Millisecond)
			return nil
		})

		if ok {
			t.Fatalf("Expected error on .TryGo() because Pool is closed")
		}
	})

	t.Run("collects errors correctly before Wait", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p := pool.New(3).WithErrors().WithContext(ctx)
		jobCount := 50
		var completed atomic.Int64
		var errored atomic.Int64

		for i := 0; i < jobCount; i++ {
			p.Go(func(c context.Context) error {
				time.Sleep(2 * time.Millisecond)
				completed.Add(1)

				if i%5 == 0 {
					errored.Add(1)
					return fmt.Errorf("err%d", i)
				}
				return nil
			})
		}

		collectedErrs := p.Collect()
		if completed.Load() != int64(jobCount) {
			t.Errorf("Jobs expected after Collect: %d, got: %d", errored.Load(), completed.Load())
		}
		if len(collectedErrs) != int(errored.Load()) {
			t.Errorf("Errors count mismatch after Collect: count: %d, collected: %d", errored.Load(), len(collectedErrs))
		}

		errored.Store(0)

		for i := 0; i < jobCount; i++ {
			p.Go(func(c context.Context) error {
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
		if completed.Load() != int64(jobCount*2) {
			t.Errorf("Jobs expected after Wait: %d, got: %d", jobCount*2, completed.Load())
		}
		if len(errs) != int(errored.Load()) {
			t.Errorf("Errors count mismatch; count: %d, collected: %d", errored.Load(), len(errs))
		}
	})
}

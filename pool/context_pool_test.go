package pool_test

import (
	"context"
	"fmt"
	"strings"
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

		err := p.CloseAndWait()

		if ctxUsed.Load() != int64(jobCount) {
			t.Errorf("Expected all jobs to receive context, got %d/%d", ctxUsed.Load(), jobCount)
		}

		if completed.Load() != int64(jobCount) {
			t.Errorf("Jobs expected: %d, got: %d", jobCount, completed.Load())
		}

		if err == nil {
			t.Fatalf("Expected aggregated error, got nil")
		}

		errStr := strings.ReplaceAll(err.Error(), "\r\n", "\n")
		errs := strings.Split(errStr, "\n")

		if len(errs) != int(errored.Load()) {
			t.Errorf("Errors count mismatch; expected: %d, collected: %d", errored.Load(), len(errs))
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

		err := p.CloseAndWait()

		if completed.Load() != int64(jobCount) {
			t.Errorf("Jobs expected: %d, got: %d", jobCount, completed.Load())
		}

		if err != nil {
			t.Errorf("Expected error to be nil, got: %v", err)
		}
	})

	t.Run("jobs see cancellation when parent context is cancelled", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		p := pool.New(4).WithErrors().WithContext(ctx)

		sawCancel := atomic.Int64{}

		p.Go(func(c context.Context) error {
			cancel()
			<-c.Done()
			sawCancel.Add(1)
			return c.Err()
		})

		err := p.CloseAndWait()

		if sawCancel.Load() != 1 {
			t.Fatalf("Expected job to observe context cancellation")
		}

		if err == nil || err.Error() != context.Canceled.Error() {
			t.Fatalf("Expected context.Canceled error, got: %v", err)
		}
	})

	t.Run("CloseAndWait triggers internal cancel", func(t *testing.T) {
		t.Parallel()
	
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
	
		p := pool.New(2).WithErrors().WithContext(ctx)
	
		seenCancel := atomic.Bool{}
	
		p.Go(func(c context.Context) error {
			<-c.Done()
			seenCancel.Store(true)
			return nil
		})
	
		err := p.CloseAndWait()
	
		if err != nil && err != context.Canceled {
			t.Errorf("Unexpected error: %v", err)
		}
	
		if !seenCancel.Load() {
			t.Fatalf("Expected context cancellation to propagate to job")
		}
	})
}

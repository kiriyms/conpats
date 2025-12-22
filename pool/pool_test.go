package pool_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/kiriyms/conpats/pool"
)

var basicCases = []struct { name string; workers int; jobCount int } {
	{"basic", 7, 50},
	{"single worker", 1, 20},
	{"many workers", 20, 100},
	{"zero workers", 0, 20},
	{"negative workers", -5, 20},
}

func TestPool(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		p := pool.New(7)
		jobCount := 50
		var completed atomic.Int64

		for i := 0; i < jobCount; i++ {
			p.Go(func() {
				time.Sleep(2 * time.Millisecond)
				completed.Add(1)
			})
		}

		p.Wait()
		if completed.Load() != int64(jobCount) {
			t.Errorf("Jobs expected: %d, got: %d", jobCount, completed.Load())
		}
	})

	t.Run("is reuseable after Collect", func(t *testing.T) {
		t.Parallel()

		p := pool.New(7)
		jobCount := 50
		var completedTotal atomic.Int64

		var completedPartA atomic.Int64
		for i := 0; i < jobCount; i++ {
			p.Go(func() {
				time.Sleep(2 * time.Millisecond)
				completedTotal.Add(1)
				completedPartA.Add(1)
			})
		}

		p.Collect()
		if completedPartA.Load() != int64(jobCount) {
			t.Errorf("Part A jobs expected: %d, got: %d", jobCount, completedPartA.Load())
		}
		if completedPartA.Load() != completedTotal.Load() {
			t.Errorf("Expected Part A equal to total. Part A: %d, Total: %d", completedPartA.Load(), completedTotal.Load())
		}

		var completedPartB atomic.Int64
		for i := 0; i < jobCount; i++ {
			p.Go(func() {
				time.Sleep(2 * time.Millisecond)
				completedTotal.Add(1)
				completedPartB.Add(1)
			})
		}

		p.Wait()
		if completedPartB.Load() != int64(jobCount) {
			t.Errorf("Part B jobs expected: %d, got: %d", jobCount, completedPartB.Load())
		}
		if completedPartB.Load() != completedPartA.Load() {
			t.Errorf("Expected Part A equal to Part B. Part A: %d, Part B: %d", completedPartA.Load(), completedPartB.Load())
		}
		if completedPartB.Load()+completedPartA.Load() != completedTotal.Load() {
			t.Errorf("Expected sum of Part A and Part B to equal to total. Parts A + B: %d, Total: %d", completedPartA.Load()+completedPartB.Load(), completedTotal.Load())
		}
	})

	t.Run("returns error after Wait", func(t *testing.T) {
		t.Parallel()

		p := pool.New(3)
		ok := p.TryGo(func() { time.Sleep(2 * time.Millisecond) })

		if !ok {
			t.Errorf("Should not error on .Go() if the Pool is not closed yet")
		}

		p.Wait()
		ok = p.TryGo(func() { time.Sleep(2 * time.Millisecond) })

		if ok {
			t.Errorf("Expected error on .Go() because Pool is closed")
		}
	})

	t.Run("does not error after Collect", func(t *testing.T) {
		t.Parallel()

		p := pool.New(3)
		ok := p.TryGo(func() { time.Sleep(2 * time.Millisecond) })

		if !ok {
			t.Errorf("Should not error on .Go() if the Pool is not closed yet")
		}

		p.Collect()
		ok = p.TryGo(func() { time.Sleep(2 * time.Millisecond) })

		if !ok {
			t.Errorf("Should not error on .Go() after a regular .Collect()")
		}

		p.Wait()
	})

	t.Run("runs correctly with zero workers input", func(t *testing.T) {
		t.Parallel()

		p := pool.New(0)
		jobCount := 20
		var completed atomic.Int64

		for i := 0; i < jobCount; i++ {
			p.Go(func() {
				time.Sleep(1 * time.Millisecond)
				completed.Add(1)
			})
		}

		p.Wait()
		if completed.Load() != int64(jobCount) {
			t.Errorf("Jobs expected: %d, got: %d", jobCount, completed.Load())
		}
	})

	t.Run("runs correctly with negative workers input", func(t *testing.T) {
		t.Parallel()

		p := pool.New(-5)
		jobCount := 20
		var completed atomic.Int64

		for i := 0; i < jobCount; i++ {
			p.Go(func() {
				time.Sleep(1 * time.Millisecond)
				completed.Add(1)
			})
		}

		p.Wait()
		if completed.Load() != int64(jobCount) {
			t.Errorf("Jobs expected: %d, got: %d", jobCount, completed.Load())
		}
	})

	t.Run("does not error when calling Go after Wait", func(t *testing.T) {
		t.Parallel()

		p := pool.New(1)
		p.Wait()

		var completed atomic.Int64
		p.Go(func() {
			completed.Add(1)
		})

		if completed.Load() != 0 {
			t.Errorf("Expected 0 completed jobs, got %d", completed.Load())
		}
	})
}

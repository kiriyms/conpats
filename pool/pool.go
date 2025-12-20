package pool

import (
	"sync"
	"sync/atomic"
)

// Pool manages a fixed number of workers executing jobs.
type Pool struct {
	workers int
	jobs  chan func()

	activeWg sync.WaitGroup
	wg       sync.WaitGroup

	once   sync.Once
	closed atomic.Bool
}

// New creates a new Pool and immediately spawns all its workers.
func New(workers int) *Pool {
	if workers <= 0 {
		workers = 1
	}

	p := &Pool{
		workers: workers,
		jobs:  make(chan func()),
	}

	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for job := range p.jobs {
				job()
			}
		}()
	}

	return p
}

// Go submits a job to the pool.
//
// If a job is submitted after Wait() has been called, it will be dropped silently.
func (p *Pool) Go(job func()) {
	closed := p.closed.Load()

	if closed {
		return
	}

	p.activeWg.Add(1)
	p.jobs <- func() {
		defer p.activeWg.Done()
		job()
	}
}

// TryGo attempts to submit a job to the pool.
//
// If a job is submitted after Wait() has been called, it will be dropped and false is returned.
// Otherwise, true is returned.
func (p *Pool) TryGo(job func()) bool {
	closed := p.closed.Load()

	if closed {
		return false
	}

	p.activeWg.Add(1)
	p.jobs <- func() {
		defer p.activeWg.Done()
		job()
	}

	return true
}

// Collect blocks until all submitted jobs are finished.
//
// This does not prevent new jobs from being submitted after using Collect().
func (p *Pool) Collect() {
	p.activeWg.Wait()
}

// Wait closes the job queue and blocks until all workers finish the jobs.
//
// After calling Wait(), the Pool is considered closed; new jobs will be dropped.
func (p *Pool) Wait() {
	p.once.Do(func() {
		p.closed.Store(true)
		close(p.jobs)
	})

	p.wg.Wait()
}

// WithErrors converts the Pool to an ErrorPool
//
// Error pool can collect errors from jobs.
func (p *Pool) WithErrors(onlyFirstErr bool) *ErrorPool {
	return &ErrorPool{
		pool:         p,
		onlyFirstErr: onlyFirstErr,
	}
}

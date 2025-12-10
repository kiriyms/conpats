package pool

import (
	"sync"
	"sync/atomic"
)

// Pool manages a fixed number of workers executing Jobs.
type Pool struct {
	limit int
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
		limit: workers,
		jobs:  make(chan func()),
	}

	for i := 0; i < p.limit; i++ {
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
// If a job is submitted after CloseAndWait() has been called, it will be dropped silently.
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

func (p *Pool) Collect() {
	p.activeWg.Wait()
}

// Wait closes the job queue and blocks until all workers finish the jobs.
func (p *Pool) Wait() {
	p.once.Do(func() {
		p.closed.Store(true)

		// Signal no more jobs.
		close(p.jobs)
	})

	// Wait for all workers to finish.
	p.wg.Wait()
}

func (p *Pool) WithErrors() *ErrorPool {
	return &ErrorPool{
		pool: p,
	}
}

package pool

import (
	"sync"
)

// Job represents a single unit of work.
type Job func()

// Pool manages a fixed number of workers executing Jobs.
type Pool struct {
	limit int
	jobs  chan Job

	activeWg sync.WaitGroup
	wg       sync.WaitGroup

	once   sync.Once
	mu     sync.Mutex
	closed bool
}

// New creates a new Pool and immediately spawns all its workers.
func New(workers int) *Pool {
	if workers <= 0 {
		workers = 1
	}

	p := &Pool{
		limit: workers,
		jobs:  make(chan Job),
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
func (p *Pool) Go(job Job) bool {
	p.mu.Lock()
	closed := p.closed
	p.mu.Unlock()

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

func (p *Pool) Wait() {
	p.activeWg.Wait()
}

// CloseAndWait closes the job queue and blocks until all workers finish the jobs.
func (p *Pool) CloseAndWait() {
	p.once.Do(func() {
		p.mu.Lock()
		p.closed = true

		// Signal no more jobs.
		close(p.jobs)
		p.mu.Unlock()
	})

	// Wait for all workers to finish.
	p.wg.Wait()
}

func (p *Pool) WithErrors() *ErrorPool {
	return &ErrorPool{
		pool: p,
	}
}

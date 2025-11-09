package pool

import (
	"errors"
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
	mu     sync.RWMutex
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
				p.activeWg.Add(1)
				job()
				p.activeWg.Done()
			}
		}()
	}

	return p
}

// Go submits a job to the pool.
//
// If a job is submitted after CloseAndWait() has been called, it will be dropped silently.
func (p *Pool) Go(job Job) error {
	p.mu.RLock()
	closed := p.closed
	p.mu.RUnlock()

	if closed {
		return errors.New("Pool already closed")
	}
	p.jobs <- job
	return nil
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

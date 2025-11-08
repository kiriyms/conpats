package pool

import "sync"

// Job represents a single unit of work.
type Job func()

// Pool manages a fixed number of workers executing Jobs.
type Pool struct {
	limit int
	jobs  chan Job
	wg    sync.WaitGroup
}

// New creates a new Pool and immediately spawns 'limit' workers.
func New(limit int) *Pool {
	if limit <= 0 {
		limit = 1
	}

	p := &Pool{
		limit: limit,
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

// Add submits a job to the pool.
func (p *Pool) Add(job Job) {
	p.jobs <- job
}

// Wait closes the job queue and blocks until all worker finish the jobs.
func (p *Pool) Wait() {
	// Signal no more jobs.
	close(p.jobs)

	// Wait for all workers to finish.
	p.wg.Wait()
}

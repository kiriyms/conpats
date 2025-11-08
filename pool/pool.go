package pool

import "sync"

// Job represents a single unit of work.
type Job func()

// Pool manages a fixed number of workers executing Jobs.
type Pool struct {
	limit  int
	jobs   chan Job
	active sync.WaitGroup
	wg     sync.WaitGroup
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
				p.active.Add(1)
				job()
				p.active.Done()
			}
		}()
	}

	return p
}

// Add submits a job to the pool.
func (p *Pool) Add(job Job) {
	p.jobs <- job
}

func (p *Pool) Wait() {
	p.active.Wait()
}

// CloseAndWait closes the job queue and blocks until all workers finish the jobs.
func (p *Pool) CloseAndWait() {
	// Signal no more jobs.
	close(p.jobs)

	// Wait for all workers to finish.
	p.wg.Wait()
}

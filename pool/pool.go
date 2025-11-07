package pool

import "sync"

// Job represents a single unit of work.
type Job func() (any, error)

// Result is the outcome of a Job.
type Result struct {
	Value any
	Err   error
}

// Pool manages a fixed number of workers executing Jobs.
type Pool struct {
	limit   int
	jobs    chan Job
	results chan Result
	wg      sync.WaitGroup

	mu     sync.Mutex
	output []Result
	done   chan struct{}
}

// New creates a new Pool and immediately spawns 'limit' workers.
func New(limit int) *Pool {
	if limit <= 0 {
		limit = 1
	}

	p := &Pool{
		limit:   limit,
		jobs:    make(chan Job),    // unbuffered
		results: make(chan Result), // unbuffered
		done:    make(chan struct{}),
	}

	// Start result collector first, so workers never block on sending results.
	go p.collectResults()

	// Spawn all workers.
	for i := 0; i < p.limit; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for job := range p.jobs {
				val, err := job()
				p.results <- Result{Value: val, Err: err}
			}
		}()
	}

	return p
}

// Add submits a job to the pool.
func (p *Pool) Add(job Job) {
	p.jobs <- job
}

// CloseAndWait closes the job queue and waits for all results to be aggregated.
func (p *Pool) CloseAndWait() []Result {
	// Signal no more jobs.
	close(p.jobs)

	// Wait for all workers to finish.
	p.wg.Wait()

	// Close results channel (no more results will arrive).
	close(p.results)

	// Wait for collector to finish aggregating.
	<-p.done

	p.mu.Lock()
	defer p.mu.Unlock()
	return p.output
}

// collectResults runs in a dedicated goroutine to aggregate results as they come in.
func (p *Pool) collectResults() {
	for r := range p.results {
		p.mu.Lock()
		p.output = append(p.output, r)
		p.mu.Unlock()
	}
	close(p.done)
}

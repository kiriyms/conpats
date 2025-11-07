package pool

import "sync"

// Pool is ...
type Pool struct {
	limit   int
	workers int
	jobs    chan Job
	results chan Result
	wg      sync.WaitGroup
	mu      sync.Mutex
}

// New ...
func New(maxGoroutines int) *Pool {
	if maxGoroutines <= 0 {
		maxGoroutines = 1
	}

	return &Pool{
		jobs:    make(chan Job),
		results: make(chan Result),
		limit:   maxGoroutines,
	}
}

// Add ...
func (p *Pool) Add(j Job) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.workers < 1 {
		p.spawnWorker()
	}

	select {
	case p.jobs <- j:
	default:
		if p.workers < p.limit {
			p.spawnWorker()
		}
		p.jobs <- j
	}
}

func (p *Pool) CloseAndWait() []Result {
	close(p.jobs)

	p.wg.Wait()
	close(p.results)

	var out []Result
	for r := range p.results {
		out = append(out, r)
	}

	return out
}

func (p *Pool) spawnWorker() {
	p.workers++
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for job := range p.jobs {
			val, err := job()
			p.results <- Result{val, err}
		}
	}()
}

type Job func() (any, error)

type Result struct {
	Value any
	Err   error
}

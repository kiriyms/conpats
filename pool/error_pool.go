package pool

import (
	"context"
	"sync"
)

// ErrorPool extends Pool to handle jobs that return errors.
//
// A new error pool must be created using New().WithErrors(). Jobs can be submitted using Go() or TryGo().
// The error pool can be gracefully shut down using Wait(), which blocks until all submitted jobs are complete and returns collected errors.
type ErrorPool struct {
	pool *Pool

	onlyFirstErr bool

	mu   sync.Mutex
	errs []error
}

// Go submits a job to the error pool.
//
// If a job is submitted after Wait() has been called, it will be dropped silently.
func (p *ErrorPool) Go(job func() error) {
	p.pool.Go(func() {
		p.addErr(job())
	})
}

// TryGo attempts to submit a job to the error pool.
//
// If a job is submitted after Wait() has been called, it will be dropped and false is returned.
// Otherwise, true is returned.
func (p *ErrorPool) TryGo(job func() error) bool {
	return p.pool.TryGo(func() {
		p.addErr(job())
	})
}

// Collect blocks until all submitted jobs are finished and returns collected errors.
//
// This does not prevent new jobs from being submitted after using Collect().
// Collect() does not close the error pool and stop the goroutine workers.
func (p *ErrorPool) Collect() []error {
	p.pool.Collect()
	return p.getErrs()
}

// Wait closes the job queue and blocks until all workers finish the jobs and returns collected errors.
//
// After calling Wait(), the error pool is considered closed; new jobs will be dropped.
func (p *ErrorPool) Wait() []error {
	p.pool.Wait()
	return p.getErrs()
}

// WithErrors converts the ErrorPool to a ContextPool
//
// ContextPool accepts jobs that expect a ctx.context as a parameter and can return errors.
func (p *ErrorPool) WithContext(ctx context.Context) *ContextPool {
	cctx, cancel := context.WithCancel(ctx)

	return &ContextPool{
		errorPool: p,
		ctx:       cctx,
		cancel:    cancel,
	}
}

func (p *ErrorPool) getErrs() []error {
	p.mu.Lock()
	errs := p.errs
	p.errs = nil
	p.mu.Unlock()

	if len(errs) == 0 {
		return nil
	}

	if p.onlyFirstErr {
		return []error{errs[0]}
	}
	return errs
}

func (p *ErrorPool) addErr(err error) {
	if err != nil {
		p.mu.Lock()
		p.errs = append(p.errs, err)
		p.mu.Unlock()
	}
}

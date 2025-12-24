package pool

import (
	"context"
)

// ContextPool extends ErrorPool to handle jobs that expect a context.Context as a parameter and can return errors.
//
// A new context pool must be created using New().WithErrors().WithContext(). Jobs can be submitted using Go() or TryGo().
// The context pool can be gracefully shut down using Wait(), which blocks until all submitted jobs are complete and returns collected errors.
type ContextPool struct {
	errorPool *ErrorPool

	cancelOnErr bool

	ctx    context.Context
	cancel context.CancelFunc
}

// Go submits a job to the context pool.
//
// If a job is submitted after Wait() has been called, it will be dropped silently.
func (p *ContextPool) Go(job func(context.Context) error) {
	p.errorPool.Go(func() error {
		err := job(p.ctx)
		if err != nil && p.cancelOnErr {
			p.cancel()
			if p.errorPool.onlyFirstErr {
				p.errorPool.addErr(err)
			}
		}

		return err
	})
}

// TryGo attempts to submit a job to the context pool.
//
// If a job is submitted after Wait() has been called, it will be dropped and false is returned.
// Otherwise, true is returned.
func (p *ContextPool) TryGo(job func(context.Context) error) bool {
	return p.errorPool.TryGo(func() error {
		err := job(p.ctx)
		if err != nil && p.cancelOnErr {
			p.cancel()
			if p.errorPool.onlyFirstErr {
				p.errorPool.addErr(err)
			}
		}

		return err
	})
}

// Collect blocks until all submitted jobs are finished and returns collected errors.
//
// This does not prevent new jobs from being submitted after using Collect().
// Collect() does not close the context pool and stop the goroutine workers.
func (p *ContextPool) Collect() []error {
	return p.errorPool.Collect()
}

// Wait closes the job queue and blocks until all workers finish the jobs and returns collected errors.
//
// After calling Wait(), the context pool is considered closed; new jobs will be dropped.
func (p *ContextPool) Wait() []error {
	if p.cancel != nil {
		defer p.cancel()
	}

	err := p.errorPool.Wait()

	return err
}

// WithCancelOnError sets whether the context pool should cancel its context upon encountering an error in any job.
// By default, this is false and all the jobs will continue to run even if some jobs return errors.
func (p *ContextPool) WithCancelOnError(cancelOnErr bool) *ContextPool {
	p.cancelOnErr = cancelOnErr
	return p
}

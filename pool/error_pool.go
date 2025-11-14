package pool

import (
	"context"
	"errors"
	"sync"
)

type ErrorJob func() error

type ErrorPool struct {
	pool *Pool

	mu   sync.Mutex
	errs []error
}

func (p *ErrorPool) Go(job ErrorJob) bool {
	return p.pool.Go(func() {
		p.addErr(job())
	})
}

func (p *ErrorPool) Wait() error {
	p.pool.Wait()
	return p.getErrs()
}

func (p *ErrorPool) CloseAndWait() error {
	p.pool.CloseAndWait()
	return p.getErrs()
}

func (p *ErrorPool) WithContext(ctx context.Context) *ContextPool {
	return &ContextPool{
		errorPool: p,
		ctx:       ctx,
	}
}

func (p *ErrorPool) getErrs() error {
	p.mu.Lock()
	errs := p.errs
	p.errs = nil
	p.mu.Unlock()

	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

func (p *ErrorPool) addErr(err error) {
	if err != nil {
		p.mu.Lock()
		p.errs = append(p.errs, err)
		p.mu.Unlock()
	}
}

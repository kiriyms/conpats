package pool

import (
	"context"
	"sync"
)

type ErrorPool struct {
	pool *Pool

	mu   sync.Mutex
	errs []error
}

func (p *ErrorPool) Go(job func() error) {
	p.pool.Go(func() {
		p.addErr(job())
	})
}

func (p *ErrorPool) TryGo(job func() error) bool {
	return p.pool.TryGo(func() {
		p.addErr(job())
	})
}

func (p *ErrorPool) Collect() []error {
	p.pool.Collect()
	return p.getErrs()
}

func (p *ErrorPool) Wait() []error {
	p.pool.Wait()
	return p.getErrs()
}

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
	return errs
}

func (p *ErrorPool) addErr(err error) {
	if err != nil {
		p.mu.Lock()
		p.errs = append(p.errs, err)
		p.mu.Unlock()
	}
}

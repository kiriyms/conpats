package pool

import (
	"context"
)

type ContextPool struct {
	errorPool *ErrorPool

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *ContextPool) Go(job func(context.Context) error) {
	p.errorPool.Go(func() error {
		return job(p.ctx)
	})
}

func (p *ContextPool) TryGo(job func(context.Context) error) bool {
	return p.errorPool.TryGo(func() error {
		return job(p.ctx)
	})
}

func (p *ContextPool) Collect() []error {
	return p.errorPool.Collect()
}

func (p *ContextPool) Wait() []error {
	if p.cancel != nil {
		p.cancel()
	}

	err := p.errorPool.Wait()

	return err
}

package pool

import (
	"context"
	"errors"
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

func (p *ContextPool) Wait() error {
	return p.errorPool.Collect()
}

func (p *ContextPool) CloseAndWait() error {
	if p.cancel != nil {
		p.cancel()
	}

	err := p.errorPool.Wait()

	if errors.Is(err, context.Canceled) {
		return context.Canceled
	}

	return err
}

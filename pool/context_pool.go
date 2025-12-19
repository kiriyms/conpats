package pool

import (
	"context"
)

type ContextPool struct {
	errorPool *ErrorPool

	cancelOnErr bool

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *ContextPool) Go(job func(context.Context) error) {
	p.errorPool.Go(func() error {
		err := job(p.ctx)
		if err != nil && p.cancelOnErr {
			p.cancel()
			p.errorPool.addErr(err)
		}

		return err
	})
}

func (p *ContextPool) TryGo(job func(context.Context) error) bool {
	return p.errorPool.TryGo(func() error {
		err := job(p.ctx)
		if err != nil && p.cancelOnErr {
			p.cancel()
			p.errorPool.addErr(err)
		}

		return err
	})
}

func (p *ContextPool) Collect() []error {
	return p.errorPool.Collect()
}

func (p *ContextPool) Wait() []error {
	if p.cancel != nil {
		defer p.cancel()
	}

	err := p.errorPool.Wait()

	return err
}

func (p *ContextPool) WithCancelOnError(cancelOnErr bool) *ContextPool {
	p.cancelOnErr = cancelOnErr
	return p
}

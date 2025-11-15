package pool

import "context"

type ContextJob func(ctx context.Context) error

type ContextPool struct {
	errorPool *ErrorPool

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *ContextPool) Go(job ContextJob) bool {
	return p.errorPool.Go(func() error {
		return job(p.ctx)
	})
}

func (p *ContextPool) Wait() error {
	return p.errorPool.Wait()
}

func (p *ContextPool) CloseAndWait() error {
	if p.cancel != nil {
		defer p.cancel()
	}
	return p.errorPool.CloseAndWait()
}

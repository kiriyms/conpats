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

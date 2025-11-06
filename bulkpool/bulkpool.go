package bulkpool

import "sync"

type BulkPool[T any, R any] struct {
	workers int
	jobs    []T
	process func(T) R
	wg      sync.WaitGroup
}

func New[T any, R any](workers int, jobs []T, process func(T) R) *BulkPool[T, R] {
	return &BulkPool[T, R]{
		workers: workers,
		jobs:    jobs,
		process: process,
	}
}

func (bp *BulkPool[T, R]) Run() []R {
	results := make([]R, len(bp.jobs))
	jobChan := make(chan int, len(bp.jobs))

	for i := 0; i < bp.workers; i++ {
		bp.wg.Add(1)
		go func() {
			defer bp.wg.Done()
			for idx := range jobChan {
				results[idx] = bp.process(bp.jobs[idx])
			}
		}()
	}

	for i := 0; i < len(bp.jobs); i++ {
		jobChan <- i
	}
	close(jobChan)

	bp.wg.Wait()
	return results
}

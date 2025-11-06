package bulkpool

type BulkPool struct {
	workers int
	jobs    []int
	process func(int) int
	wg      sync.WaitGroup
}

func New(workers int, jobs []int, process func(int) int) *BulkPool {
	return &BulkPool{
		workers: workers,
		jobs:    jobs,
		process: process,
	}
}

func (bp *BulkPool) Run() []int {
	results := make([]int, len(bp.jobs))
	jobChan := make(chan int, len(bp.jobs))
	for i := 0; i < bp.workers; i++ {
		bp.wg.Add(1)
		go func() {
			defer bp.wg.Done()
			for job := range jobChan {
				results[job] = bp.process(job)
			}
		}()
	}

	for i := 0; i < len(bp.jobs); i++ {
		jobChan <- bp.jobs[i]
	}
	close(jobChan)
	bp.wg.Wait()
	return results
}

package conpats

import (
	"fmt"
	"sync"
)

func WorkerPool(numWorkers int, jobs []int, f func(int) int) <-chan int {
	in := make(chan int, len(jobs))
	out := make(chan int, len(jobs))
	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, in, out, f, &wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func worker(id int, in <-chan int, out chan<- int, f func(int) int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range in {
		fmt.Printf("Worker %d processing job %d\n", id, j)
		out <- f(j)
	}
}

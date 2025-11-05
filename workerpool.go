package conpats

import (
	"fmt"
	"sync"
)

func Worker(id int, in <-chan int, out chan<- int, f func(int) int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range in {
		fmt.Printf("Worker %d processing job %d\n", id, j)
		out <- f(j)
	}
}

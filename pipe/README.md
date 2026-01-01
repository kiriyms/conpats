## Pipeline

Pipeline API implements a concurrency pattern, where an input of `values[I]` is passed through a function `func(I) O`, transforming them to an output of `values[O]`.

### Usage

The **Pipes** are constructed using one of the two functions: [`pipe.PipeFromChan(...)`](https://pkg.go.dev/github.com/kiriyms/conpats/pipe#PipeFromChan):

```go
nums := []int{1, 2, 3, 4, 5}

in := make(chan int)
go func() {
    defer close(in)
    for _, n := range nums {
        in <- n
    }
}()

out := pipe.PipeFromChan(func(n int) int {
    return n * n
}, in, 2)
```

Or [`pipe.PipeFromSlice(...)`](https://pkg.go.dev/github.com/kiriyms/conpats/pipe#PipeFromSlice), which works like the above example, but abstracts the generator:

```go
nums := []int{1, 2, 3, 4, 5}

out := pipe.PipeFromSlice(func(n int) int {
    return n * n
}, nums, 2)
```

**Pipes** always return a `channel[O]` of values, making it easy to chain several pipes together, creating a full **Pipeline**:

```go
nums := []int{1, 2, 3, 4, 5}

squareCh := pipe.PipeFromSlice(func(n int) int {
    return n * n
}, nums, 2)

incrementCh := pipe.PipeFromChan(func(n int) int {
    return n + 1
}, squareCh, 5)

doubleCh := pipe.PipeFromChan(func(n int) int {
    return n * 2
}, incrementCh, 10)
```

**Pipes** don't have to adhere to a single common type. Easily chain several **Pipes** transforming from one type to another:

```go
nums := []int{1, 2, 3, 4, 5}

sqrtCh := pipe.PipeFromSlice(func(n int) float64 {
    return math.Sqrt(float64(n))
}, nums, 3)

logChan := pipe.PipeFromChan(func(n float64) string {
    return fmt.Sprintf("Sqrt: %.2f", n)
}, sqrtChan, 1)
```

Conveniently collect the results of a final **Pipe** segment using a utility function [`pipe.Collect(chan)`](https://pkg.go.dev/github.com/kiriyms/conpats/pipe#Collect), which will block until the **Pipe** output channel is closed:

```go
nums := []int{1, 2, 3, 4, 5}

out := pipe.PipeFromSlice(func(n int) float64 {
    return math.Sqrt(float64(n))
}, nums, 3)

results := pipe.Collect(out)
```

The **Pipes** process values concurrently using a **Worker Pool** under the hood. By default, [`pool.Pool`](https://pkg.go.dev/github.com/kiriyms/conpats/pool#Pool) provided by `conpats` is used. However, the **Pool** implementation can be configured using the [`pipe.WithPool(pool)`](https://pkg.go.dev/github.com/kiriyms/conpats/pipe#WithPool), which accepts a simple [`pipe.Pool`](https://pkg.go.dev/github.com/kiriyms/conpats/pipe#Pool) interface:

```go
type Pool interface {
	Go(func())
	Wait()
}
```

```go
customPool := NewCustomPool()

nums := []int{1, 2, 3, 4, 5}

out := pipe.PipeFromSlice(func(n int) float64 {
    return math.Sqrt(float64(n))
}, nums, 3, WithPool(customPool))

results := pipe.Collect(out)
```

To avoid goroutine leaks, it is expected that the `Wait()` function of the [`pipe.Pool`](https://pkg.go.dev/github.com/kiriyms/conpats/pipe#Pool) interface is the final step of the **Pool** lifecycle, which closes the **Pool** and blocks until all remaining work is done and the underlying goroutines are terminated.

Alternatively, provide a **Worker Pool** implementation from a different package. For example, [`conc`'s Pool](https://github.com/sourcegraph/conc) fits the interface:

```go
import "github.com/sourcegraph/conc/pool"

func main() {
    concPool := pool.New()

    nums := []int{1, 2, 3, 4, 5}

    out := pipe.PipeFromSlice(func(n int) float64 {
        return math.Sqrt(float64(n))
    }, nums, 3, WithPool(concPool))

    results := pipe.Collect(out)
}
```

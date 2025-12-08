# `conpats`: collection of structured concurrency patterns for Go

`conpats` contains several common concurrency patterns for convenient use.

## Quick rundown

- Use `pool.Pool` if you want a reusable Worker Pool.
- Use `pool.ErrorPool` if your Worker Pool tasks return errors.
- Use `pool.ContextPool` if your Worker Pool tasks use context.

- Use `pipeline.Pipeline` if you want to process several functions in sequence using concurrency

## Goals

Main goals of this package are:

1. Make concurrency easier and reduce boilerplate
2. Provide a variety of common concurrency patterns in one place
3. Structurize concurrency to improve control
4. Avoid any third-party dependencies

## Workflow

### Pool

- Both Pools start with `pool.New()`.
- To convert a regular Pool to an Error Pool, use `pool.New().WithErrors()`.
- To convert an Error Pool to a Context Pool, use `pool.New().WithErrors().WithContext()`.
- Add new jobs to your Pools using `p.Go()`.
- Use `p.Wait()` to block until all given jobs finish. After `p.Wait()`, you can reuse the Pool and add jobs again.
- Use `p.CloseAndWait()` to close the jobs channel and block until all given jobs finish.

### Example

Pool:

```go
func main() {
	p := pool.New(5)
	for i := range 50 {
		p.Go(func() {
			fmt.Printf("Hello - %d\n", i)
		})
	}
	p.Wait()
	for i := range 10 {
		p.Go(func() {
			fmt.Printf("Bye - %d\n", i)
		})
	}
	p.CloseAndWait()
}
```

Error Pool:

```go
func main() {
    p := pool.New(7).WithErrors()
    jobCount := 50

    for i := 0; i < 50; i++ {
        p.Go(func() error {
            time.Sleep(2 * time.Millisecond)
            return nil
        })
    }

    err := p.Wait()
    fmt.Println(err) // err == nil

    for i := 0; i < jobCount; i++ {
        p.Go(func() error {
            time.Sleep(2 * time.Millisecond)

            if i%5 == 0 {
                errored.Add(1)
                return fmt.Errorf("err-%d", i)
            }
            return nil
        })
    }

    err = p.CloseAndWait()
    fmt.Println(err) // err == "err-0 err-5 err-10 ..."
}
```

### Pipeline

- Start by creating a new Pipeline from either a slice `pipeline.NewFromSlice()` or a channel `pipeline.NewFromChannel()`.
- Set up the stages of the pipeline by passing you own functions to `pipeline.AddStage()`.
- Run the pipeline using `pipeline.Run()`. This returns a channel of final output values, which you can read or pass to another pipeline with `pipeline.NewFromChannel()`.

### Example

Pipeline:

```go
func main() {
	pipe := NewFromSlice([]int{1, 2, 3, 4, 5})
	pipe.AddStage(func(i int) int {
		return i * 2
	})
	pipe.AddStage(func(i int) int {
		return i + 1
	})
	out := pipe.Run()

	res := make([]int, 0)
	for v := range out {
		res = append(res, v)
	}

	fmt.Println(res) // res == [3 5 7 9 11]
}
```

## Characteristics

- Worker goroutines in a Pool are created upfront, when `pool.New()` is called.
- `pool.Go()` blocks if all workers are busy.
- `pool.Go()` returns `false`, then drops the job if called after the Pool has been closed.
- Error Pool returns all accumulated errors as a single `error` on `p.Wait()` and `p.CloseAndWait()`.
- If no jobs errored, Error Pool will return `nil`.

## Status

This package is in a `0.1.0` version.

Core work is in progress.
Macro-goals are:

1. Flesh-out Pool and Error Pool, add Lazy Pool
2. Add more concurrency patters, such as Fan-Out/Fan-In, Pub-Sub and others.
3. Add monadic concurrency support for Pools

### TODO

1. Examples - Load Balancer utilizing WorkerPool
2. Make patterns modular - one using another
3. If possible, make patterns compatible with implementations from other packages
4. Add "Thoughts&Notes"
5. Add cookbook

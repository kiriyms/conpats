# `conpats`

`conpats` contains several common concurrency patterns for convenient use.

```sh
go get github.com/kiriyms/conpats
```

## Quick rundown

- `conpats` provides **Worker Pool**, **Pipeline** and **Tee**.

#### [Worker pool](/pool/README.md)

- Use [`pool.Pool`](/pool/pool.go) when you need to run jobs concurrently with a goroutine limit.
- Use [`pool.ErrorPool`](/pool/error_pool.go) when you need to run jobs _that return errors_ concurrently with a giroutine limit.
- Use [`pool.ContextPool`](/pool/context_pool.go) when you need to run jobs _that return errors and receive a `ctx.Context` parameter_ concurrently with a giroutine limit.

Every **Pool** must be created using [`pool.New()`](/pool/pool.go). To convert it use:

- `.New().WithError(bool)` to get a `pool.ErrorPool`, where the `bool` parameter specifies if you want only the first error to be returned, rather that a slice of all errors.
- `.New().WithError(bool).WithContext(ctx)` to get a `pool.ContextPool`, where the `ctx` paramater specifies your parent context that needs to be passed to all your jobs.

#### [Pipeline](/pipe/README.md)

- Use [`pipe.PipeFromChan(...)`](/pipe/pipe.go) when you need to run all input values from a given channel through a function concurrently.
- Use [`pipe.PipeFromSlice(...)`](/pipe/pipe.go) when you need to run all values of a given slice through a function concurrently.

Both **Pipe** functions return channels, making it easy to chain several pipes together or using the output channel in other ways, for example:

- Use [`pipe.Collect(chan)`](/pipe/pipe.go) when you want to block and collect results from a channel into a slice until it is closed.

The **Pipeline** implementation uses the [`pool.Pool`](/pool/pool.go) by default, but can be modified:

- Use [`pipe.WithPool(pool)`](/pipe/pipe.go) option parameter to specify the **Worker Pool** implementation that the **Pipe** will use.

#### [Tee](/tee/README.md)

- Use [`tee.NewTee(chan)`](/tee/tee.go) to create several channels (buffered or unbuffered) that each receive a copy of a value from a provided `chan` channel.

## Goals

Main goals of this package are:

1. Make concurrency easier and reduce boilerplate
2. Provide a variety of common concurrency patterns in one place
3. Avoid any third-party dependencies

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

This package is in a `1.0.0` version.

Core work is done, common concurrency patterns are implemented.
Possible future improvements:

- Add more patters & utility functions (like Fan-in/Fan-out, Pub-Sub, etc.)

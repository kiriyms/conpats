# `conpats`: collection of structured concurrency patterns for Go

`conpats` contains several common concurrency patterns for convenient use.

## Quick rundown

- Use `pool.Pool` if you want a reusable Worker Pool.
- Use `pool.ErrorPool` if your Worker Pool tasks return errors.

## Goals

Main goals of this package are:

1. Make concurrency easier and reduce boilerplate
2. Provide a variety of common concurrency patterns in one place
3. Structurize concurrency to improve control
4. Avoid any third-party dependencies

## Workflow

- Both Pools start with `pool.New()`.
- To convert a regular Pool to an Error Pool, use `pool.New().WithErrors()`.
- Add new jobs to your Pools using `p.Go()`.
- Use `p.Wait()` to block until all given jobs finish. After `p.Wait()`, you can reuse the Pool and add jobs again.
- Use `p.CloseAndWait()` to close the jobs channel and block until all given jobs finish.

### Example

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

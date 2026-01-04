![conpats banner](/conpats_banner.svg)

# `conpats`

`conpats` contains several common concurrency patterns for convenient use.

<div align="center" style="display: flex; justify-content: center; gap: 16px;">

<div>

[![Go Reference](https://pkg.go.dev/badge/github.com/kiriyms/conpats.svg)](https://pkg.go.dev/github.com/kiriyms/conpats)

</div>

<div>

[![Go Report Card](https://goreportcard.com/badge/github.com/kiriyms/conpats)](https://goreportcard.com/report/github.com/kiriyms/conpats)

</div>

<div>

[![codecov](https://codecov.io/gh/kiriyms/conpats/branch/main/graph/badge.svg)](https://codecov.io/gh/kiriyms/conpats)

</div>

<div>

[![Tag](https://img.shields.io/github/v/tag/kiriyms/conpats?style=flat-square&logo=fitbit&color=%23ff8936)](https://github.com/kiriyms/conpats/tags)

</div>

</div>

```sh
go get github.com/kiriyms/conpats
```

## Table of Contents

- [Quick Rundown](#quick-rundown)
  - [Worker Pool](#worker-pool)
  - [Pipeline](#pipeline)
  - [Tee](#tee)
- [Goals](#goals)
- [Usage](#usage)
  - [Worker Pool](#worker-pool-1)
  - [Pipeline](#pipeline-1)
  - [Tee](#tee-1)
- [Cookbook](#cookbook)
- [Thoughts & Notes](#thoughts--notes)
- [Status](#status)

## Quick Rundown

- `conpats` provides **Worker Pool**, **Pipeline** and **Tee**.

#### [Worker Pool](/pool/README.md)

- Use [`pool.Pool`](https://pkg.go.dev/github.com/kiriyms/conpats/pool#Pool) when you need to run jobs concurrently with a goroutine limit.
- Use [`pool.ErrorPool`](https://pkg.go.dev/github.com/kiriyms/conpats/pool#ErrorPool) when you need to run jobs _that return errors_ concurrently with a giroutine limit.
- Use [`pool.ContextPool`](https://pkg.go.dev/github.com/kiriyms/conpats/pool#ContextPool) when you need to run jobs _that return errors and receive a `context.Context` argument_ concurrently with a giroutine limit.

Every **Pool** must be created using [`pool.New(...)`](https://pkg.go.dev/github.com/kiriyms/conpats/pool#New). To convert it use:

- `.New(...).WithErrors()` to get a `pool.ErrorPool`.
- `.New(...).WithErrors().WithContext(ctx)` to get a `pool.ContextPool`, where the `ctx` paramater specifies your parent context that needs to be passed to all your jobs.

#### [Pipeline](/pipe/README.md)

- Use [`pipe.PipeFromChan(...)`](https://pkg.go.dev/github.com/kiriyms/conpats/pipe#PipeFromChan) when you need to run all input values from a given channel through a function concurrently.
- Use [`pipe.PipeFromSlice(...)`](https://pkg.go.dev/github.com/kiriyms/conpats/pipe#PipeFromSlice) when you need to run all values of a given slice through a function concurrently.

Both **Pipe** functions return channels, making it easy to chain several pipes together or using the output channel in other ways, for example:

- Use [`pipe.Collect(chan)`](https://pkg.go.dev/github.com/kiriyms/conpats/pipe#Collect) when you want to block and collect results from a channel into a slice until it is closed.

The **Pipeline** implementation uses the [`pool.Pool`](https://pkg.go.dev/github.com/kiriyms/conpats/pool#Pool) by default, but can be modified:

- Use [`pipe.WithPool(pool)`](https://pkg.go.dev/github.com/kiriyms/conpats/pipe#WithPool) option parameter to specify the **Worker Pool** implementation that the **Pipe** will use.

#### [Tee](/tee/README.md)

- Use [`tee.NewTee(chan)`](https://pkg.go.dev/github.com/kiriyms/conpats/tee#NewTee) to create several channels (buffered or unbuffered) that each receive a copy of a value from a provided `chan` channel.

## Goals

Main goals of this package are:

1. Make concurrency easier and reduce boilerplate
2. Provide a variety of common concurrency patterns in one place
3. Avoid any third-party dependencies

## Usage

This section provides simple usage examples of **Worker Pool**, **Pipeline** and **Tee** usage compared to manual implementation. More examples can be found in these patterns' respective READMEs: [Pool](/pool/README.md), [Pipe](/pipe/README.md), [Tee](/tee/README.md).

#### [Worker Pool](/pool/README.md)

<table>
<thead>
<tr>
<th>Manual</th>
<th>Using <a href="/pool/pool.go"><code>pool.Pool</code></a></th>
</tr>
</thead>
<tbody>
<tr>
<td>

```go
func main() {
	wg := sync.WaitGroup{}
	jobs := make(chan func())
	for i := 0; i < 10; i++ {
		wg.Go(func() {
			for job := range jobs {
				job()
			}
		})
	}

	for i := 0; i < 100; i++ {
		jobs <- doWork
	}
	close(jobs)
	wg.Wait()
}
```

</td>
<td>

```go
func main() {
	p := pool.New(10)
	for i := 0; i < 100; i++ {
		p.Go(doWork)
	}
	p.Wait()
}
```

</td>
</tr>
</tbody>
</table>

#### [Pipeline](/pipe/README.md)

<table>
<thead>
<tr>
<th>Manual</th>
<th>Using <a href="/pipe/pipe.go"><code>pipe.PipeFromChan()</code></a></th>
</tr>
</thead>
<tbody>
<tr>
<td>

```go
func main() {
		nums := []int{1, 2, 3, 4, 5}

	in := make(chan int)
	go func() {
		defer close(in)
		for _, n := range nums {
			in <- n
		}
	}()

	sqrtChan := make(chan float64)
	wgSqrt := sync.WaitGroup{}
	go func() {
		defer close(sqrtChan)
		defer wgSqrt.Wait()
		for i := 0; i < 5; i++ {
			wgSqrt.Add(1)
			go func() {
				defer wgSqrt.Done()
				for n := range in {
					sqrtChan <- float64(math.Sqrt(float64(n)))
				}
			}()
		}
	}()

	logChan := make(chan string)
	wgLog := sync.WaitGroup{}
	go func() {
		defer close(logChan)
		defer wgLog.Wait()
		for i := 0; i < 3; i++ {
			wgLog.Add(1)
			go func() {
				defer wgLog.Done()
				for sq := range sqrtChan {
					logChan <- fmt.Sprintf("Sqrt: %.2f", sq)
				}
			}()
		}
	}()

	results := make([]string, 0)
	for log := range logChan {
		results = append(results, log)
	}
}
```

</td>
<td>

```go
func main() {
	nums := []int{1, 2, 3, 4, 5}

	sqrtChan := pipe.PipeFromSlice(func(n int) float64 {
		return math.Sqrt(float64(n))
	}, nums, 5)

	logChan := pipe.PipeFromChan(func(n float64) string {
		return fmt.Sprintf("Sqrt: %.2f", n)
	}, sqrtChan, 2)

	results := pipe.Collect(logChan)
}
```

</td>
</tr>
</tbody>
</table>

#### [Tee](/tee/README.md)

<table>
<thead>
<tr>
<th>Manual</th>
<th>Using <a href="/tee/tee.go"><code>tee.NewTee()</code></a></th>
</tr>
</thead>
<tbody>
<tr>
<td>

```go
func main() {
	in := make(chan int)
	outs := make([]chan int, 3)
	for i := range 3 {
		outs[i] = make(chan int)
	}

	go func() {
		defer func() {
			for _, out := range outs {
				close(out)
			}
		}()

		for item := range in {
			for _, out := range outs {
				out <- item
			}
		}
	}()
}
```

</td>
<td>

```go
func main() {
	in := make(chan int)
	outs := tee.NewTee(in, 3, 0)
}
```

</td>
</tr>
</tbody>
</table>

> **Note**: if one of the output channels is blocked and waiting to be read from, it will cause all other output channels to block too.

## Cookbook

The concurrency pattern abstractions in `conpats` can be easily combined with each other.

To see usage examples that are more complex and closer to real-world problems, check out the [Cookbook](/cookbook.md).

## Thoughts & Notes

Making a small `Go` package has been an enlightening and interesting experience. As a result of this endeavor, I've jotted down some [final thoughts](/thoughts.md).

## Status

**`v1`** (core API settled).

Common concurrency patterns are implemented.
Possible future improvements:

- Add more patters & utility functions (like Fan-in/Fan-out, Pub-Sub, etc.)
- Add more cookbook examples

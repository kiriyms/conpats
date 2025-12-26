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


## Status

This package is in a `1.0.0` version.

Core work is done, common concurrency patterns are implemented.
Possible future improvements:

- Add more patters & utility functions (like Fan-in/Fan-out, Pub-Sub, etc.)

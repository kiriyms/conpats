## Pool

Worker Pool API implements a concurrency pattern, where `n` goroutines are created and re-used to complete a load of jobs. The **Pool** enables _controlled_, _concurrent_ processing of jobs.

In addition to the regular [`pool.Pool`](/pool/pool.go), `conpats` provides extended **Pools**: [`pool.ErrorPool`](/pool/error_pool.go) and [`pool.ContextPool`](/pool/context_pool.go).

### Usage

The core workflow with the **Pool** is in three steps:

1. `pool.New(...)`: start by creating a new pool.
2. `pool.Go(...)`: add jobs to the pool to e processed concurrently.
3. `pool.Wait()`: block until all jobs are finished and close the pool, freeing up resources.

Every **Pool** always begins with a `pool.New(workers)` constructor, specifying the number of `worker` goroutines that will be created and used in the pool:

```go
// create a pool with 10 worker goroutines
p := pool.New(10)

// negative or zero workers are valid, but will return a pool with 1 worker
p = pool.New(0)
p = pool.New(-5)
```

> **Note**: all worker goroutines in a **Pool** are created _immediately_, and are _re-used_ for jobs.

Add an arbitrary number of jobs using `pool.Go(func())`:

```go
p := pool.New(10)

for i := 0; i < 100; i++ {
    p.Go(func() {
        fmt.Println("Processing job: ", num)
        // work
    })
}
```

> **Note**: `pool.Go(...)` will block if all worker goroutines are busy at the moment.

Finally, to syncronize and wait for all submitted work to finish, use `pool.Wait()`, which will also close the **Pool** and its workers, freeing up resources:

```go
p := pool.New(10)

for i := 0; i < 100; i++ {
    p.Go(func() {
        fmt.Println("Processing job: ", num)
        // work
    })
}

p.Wait()
```

The core idea of a **Worker Pool** is to re-use its workers. Syncronize and wait for job completion without closing the **Pool** using `pool.Collect()`, which enables adding more jobs later, if needed:

```go
p := pool.New(10)

for i := 0; i < 100; i++ {
    p.Go(func() {
        fmt.Println("Processing job: ", num)
        // work
    })
}

p.Collect() // block and wait for jobs to finish without closing the pool

for i := 0; i < 25; i++ {
    p.Go(func() {
        fmt.Println("Some additional job: ", num)
        // work
    })
}

p.Wait() // block and wait, then close the pool
```

If any new jobs are submitted to the **Pool** after it's been closed, they will be dropped silently:

```go
p.Wait()

p.Go(func() {
    fmt.Println("Last job!")
}) // nothing happens
```

To check whether the job has been successfully submitted, use `pool.TryGo(...)`:

```go
p := pool.New(1)

ok := p.TryGo(func() {
    fmt.Println("Processing job")
}) // true

p.Collect()

ok = p.TryGo(func() {
    fmt.Println("Processing job")
}) // true

p.Wait()

ok = p.TryGo(func() {
    fmt.Println("Processing job")
}) // false
```

#### Error Pool

To process jobs that return errors, use `pool.ErrorPool`:

```go
p := pool.New(10).WithErrors(false)
```

Use the `onlyFirstErr` `bool` argument in `pool.WithErrors(...)` to specify:

- `false`: make `.Collect()` and `.Wait()` return all collected errors
- `true`: make `.Collect()` and `.Wait()` return only the first error

```go
p := pool.New(2).WithErrors(false)

for i := 0; i < 50; i++ {
    p.Go(func() error {
        if i%5 == 0 {
            return fmt.Errorf("err%d", i)
        }
        return nil
    })
}

errs := p.Wait() // slice of 10 errors
```

> **Note**: currently `pool.ErrorPool` does not handle panics in any way.

Like in `pool.Pool`, use `.Collect()` to block and wait for submitted jobs to finish, without closing the **Error Pool** and return the collected errors. This will also clear the **Error Pool's** error storage, meaning all subsequent `.Collect()` and `.Wait()` calls will only return the new errors:

```go
p := pool.New(10).WithErrors(false)

for i := 0; i < 100; i++ {
    p.Go(func() error {
        // work, possibly return err
        return nil
    })
}

errs := p.Collect() // any collected errors from the previous 100 submitted jobs

for i := 0; i < 25; i++ {
    p.Go(func() error {
        // work, possibly return err
        return nil
    })
}

errs = p.Wait() // any collected errors from the previous 25 newly submitted jobs
```

#### Context Pool

To process jobs that return errors and accept a `context.Context` argument, use `pool.ContextPool`:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

p := pool.New(4).WithErrors(false).WithContext(ctx)
```

The **Context Pool** creates its own child context based on the context passed in the constructor.

The **Context Pool** can be configured to cancel its context immediately when an error is returned from the jobs using `.WithCancelOnError(bool)`:

```go
p := pool.New(12).WithErrors(false).WithContext(ctx).WithCancelOnError(true)
```

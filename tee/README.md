## Tee

Tee API implements a concurrency pattern, where an input from a singular channel is _copied_ into two channels. In `conpats`, **Tee** is generalized and copies values into an `n` number of channels (buffered or unbuffered).

### Usage

Simply use [`tee.NewTee(...)`](https://pkg.go.dev/github.com/kiriyms/conpats/tee#NewTee):

```go
in := make(chan int)

// create 3 unbuffered channels
outs := tee.NewTee(in, 3, 0)
```

Returned channels can be unbuffered or buffered. Specify that using the last argument in the [`tee.NewTee(...)`](https://pkg.go.dev/github.com/kiriyms/conpats/tee#NewTee) function:

```go
in := make(chan int)

// create 5 buffered channels with buffer size 10
outs := tee.NewTee(in, 5, 10)
```

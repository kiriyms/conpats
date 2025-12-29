## Cookbook

The concurrency pattern abstractions in `conpats` can be easily combined with each other.

In this **Cookbook** you will find some more complex examples than those in the READMEs.

### Table of Contents

- [Concurrent User Fetch](#example-1-concurrent-user-fetch)

#### Example 1. Concurrent User Fetch

```go
package main

import (
	"fmt"
	"github.com/kiriyms/conpats/pipe"
	"github.com/kiriyms/conpats/pool"
	"github.com/kiriyms/conpats/tee"
	"log/slog"
	"sync"
)

type APIUser struct {
	ID      int
	email   string
	active  bool
	address string
}

type DBUser struct {
	ID        int
	email     string
	active    bool
	address   string
	createdAt string
	updatedAt string
}

type LogUser struct {
	ID     int
	email  string
	active bool
}

func main() {
	// Initial setup before starting the work

	// Sample list of user IDs to process
	ids := []int{1, 24, 13, 55, 67, 89, 43, 21, 90, 11}
	apiResultsCh := make(chan *APIUser)

	// Split the API results into 2 separate channels
	resultsSlice := tee.NewTee(apiResultsCh, 2, 0)

	// Process each channel with its own pipe

	// Set up a database pipeline
	filterPipe := pipe.PipeFromChan(func(u *APIUser) *APIUser {
		if !u.active {
			return nil
		}
		return u
	}, resultsSlice[0], 2)

	mapPipe := pipe.PipeFromChan(func(u *APIUser) *DBUser {
		if u == nil {
			return nil
		}
		return &DBUser{
			ID:        u.ID,
			email:     u.email,
			active:    u.active,
			address:   u.address,
			createdAt: "2024-01-01",
			updatedAt: "2024-01-01",
		}
	}, filterPipe, 2)

	insertPipe := pipe.PipeFromChan(func(u *DBUser) error {
		if u == nil {
			return nil
		}
		err := insertToDatabase(u)
		return err
	}, mapPipe, 2)

	// Set up logging pipeline
	redactPipe := pipe.PipeFromChan(func(u *APIUser) *LogUser {
		return &LogUser{
			ID:     u.ID,
			email:  u.email,
			active: u.active,
		}
	}, resultsSlice[1], 2)

	logPipe := pipe.PipeFromChan(func(u *LogUser) error {
		slog.Info("User Fetched", slog.String("User", fmt.Sprintf("%+v", u)))
		return nil
	}, redactPipe, 2)

	// Set up pipeline results collection
	dbResults := make([]error, 0)
	logResults := make([]error, 0)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for res := range insertPipe {
			dbResults = append(dbResults, res)
		}
		wg.Done()
	}()

	go func() {
		for res := range logPipe {
			logResults = append(logResults, res)
		}
		wg.Done()
	}()

	// Start the work by creating a pool and submitting jobs
	p := pool.New(4)
	for _, v := range ids {
		p.Go(func() {
			apiResultsCh <- fetchUserFromAPI(v)
		})
	}

	// Wait for all jobs to complete and clean up
	p.Wait()
	close(apiResultsCh)

	// Wait for result collection to complete
	wg.Wait()

	// Use `dbResults` and `logResults` as needed
	fmt.Println("Results from db pipe:", dbResults)
	fmt.Println("Results from log pipe:", logResults)
}
```

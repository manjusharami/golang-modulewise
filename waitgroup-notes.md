# sync.WaitGroup — Detailed Notes & Practical Use

## What is it?

`WaitGroup` (from the `sync` package) is used to **wait for a collection of goroutines to finish** before continuing execution. Think of it as a counter that tracks "how many tasks are still running."

---

## How it works

| Method | Purpose |
|---|---|
| `wg.Add(n)` | Increases counter by `n` (number of goroutines you're about to launch) |
| `wg.Done()` | Decreases counter by 1 (call when a goroutine finishes, usually via `defer`) |
| `wg.Wait()` | Blocks the calling goroutine until counter reaches 0 |

### Key Rules

- Call `Add()` **before** starting the goroutine (not inside it) — avoids race conditions.
- Call `Done()` with `defer` so it always runs, even if the goroutine panics.
- Never copy a `WaitGroup` after first use — always pass by pointer (`*sync.WaitGroup`).
- Counter must never go negative (extra `Done()` calls cause a panic).
- Zero value is ready to use: `var wg sync.WaitGroup`

---

## Basic Live Example

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("All workers completed")
}
```

### Common Mistake (Race Condition)

```go
// ❌ BAD — Add() called inside the goroutine
go func() {
    wg.Add(1)   // too late, Wait() may already have returned
    defer wg.Done()
    doWork()
}()
```

Always `Add()` on the calling goroutine, **before** the `go` statement.

---

## Where WaitGroup is Used in Production

### 1. Fan-out API calls (aggregating results from multiple services)

Real scenario: an e-commerce backend needs product details, pricing, and inventory from 3 different microservices — fetch them concurrently instead of sequentially.

```go
func getProductPage(id string) ProductPage {
	var wg sync.WaitGroup
	var details Details
	var pricing Pricing
	var inventory Inventory

	wg.Add(3)

	go func() {
		defer wg.Done()
		details = fetchDetails(id)
	}()

	go func() {
		defer wg.Done()
		pricing = fetchPricing(id)
	}()

	go func() {
		defer wg.Done()
		inventory = fetchInventory(id)
	}()

	wg.Wait() // wait for all 3 API calls to finish
	return ProductPage{details, pricing, inventory}
}
```

**Why it matters:** Cuts response time from `sum(all calls)` to `max(slowest call)`.

---

### 2. Batch processing / worker pools

Real scenario: processing thousands of uploaded images/files concurrently with a limited number of workers.

```go
func processFiles(files []string) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // limit to 10 concurrent workers

	for _, file := range files {
		wg.Add(1)
		sem <- struct{}{} // acquire slot

		go func(f string) {
			defer wg.Done()
			defer func() { <-sem }() // release slot

			processFile(f)
		}(file)
	}

	wg.Wait()
	fmt.Println("All files processed")
}
```

**Why it matters:** Used in ETL pipelines, image/video processing services, log ingestion — prevents overwhelming CPU/memory or downstream systems.

---

### 3. Graceful shutdown of servers

Real scenario: an HTTP server needs to finish in-flight requests before shutting down.

```go
func main() {
	var wg sync.WaitGroup
	requests := make(chan Request)

	go func() {
		for req := range requests {
			wg.Add(1)
			go func(r Request) {
				defer wg.Done()
				handleRequest(r)
			}(req)
		}
	}()

	// on shutdown signal:
	close(requests)
	wg.Wait() // wait for in-flight requests to complete
	fmt.Println("Server shut down gracefully")
}
```

**Why it matters:** Prevents dropped requests during deploys/restarts (common in Kubernetes rolling updates).

---

### 4. Fan-in with result channel (WaitGroup + channel combo)

Real scenario: aggregating logs/metrics from multiple sources, then closing the channel once all producers finish.

```go
func mergeResults(sources []Source) <-chan Result {
	out := make(chan Result)
	var wg sync.WaitGroup

	for _, src := range sources {
		wg.Add(1)
		go func(s Source) {
			defer wg.Done()
			for r := range s.Results() {
				out <- r
			}
		}(src)
	}

	// closer goroutine: closes 'out' once all producers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
```

**Why it matters:** This pattern is common in log aggregators, metrics collectors, and stream-processing pipelines (e.g., collecting from multiple Kafka partitions concurrently).

---

## WaitGroup vs Channels — When to Use What

| Use Case | Tool |
|---|---|
| Just wait for goroutines to finish, no data returned | `WaitGroup` |
| Need to pass/collect data from goroutines | Channels |
| Need both: wait + collect results | WaitGroup + channel (see pattern above) |
| Need to limit concurrency | WaitGroup + buffered channel (semaphore pattern) |

---

## Summary Cheat Sheet

```go
var wg sync.WaitGroup

wg.Add(1)         // before starting goroutine
go func() {
    defer wg.Done()  // when goroutine finishes
    // work here
}()

wg.Wait()         // blocks until all Done() calls complete
```

**Production use cases recap:**
- Concurrent API fan-out/aggregation
- Worker pools for batch/file processing
- Graceful server shutdown
- Fan-in pipelines (logs, metrics, streams)

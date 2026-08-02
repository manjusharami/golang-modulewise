# Go Goroutines & Concurrency — Complete Notes + Practice Q&A

---

## 1. Goroutines — Basics

A **goroutine** is a lightweight thread managed by the Go runtime (not the OS). Starting one costs ~2KB of stack (grows/shrinks dynamically), so you can run thousands/millions of them.

```go
func sayHello() {
    fmt.Println("Hello")
}

func main() {
    go sayHello()        // starts a new goroutine
    time.Sleep(time.Second) // needed or main may exit before goroutine runs
}
```

**Key facts:**
- `go func(){...}()` schedules a function to run concurrently.
- The `main()` goroutine does **not** wait for other goroutines automatically. If `main` returns, the program exits — even if other goroutines haven't finished.
- Goroutines are multiplexed onto OS threads by the Go scheduler (M:N scheduling — M goroutines on N OS threads).
- `runtime.GOMAXPROCS(n)` controls how many OS threads can execute Go code simultaneously (defaults to number of CPUs).

**Common mistake — loop variable capture (pre Go 1.22):**
```go
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i) // BUG in Go <1.22: all may print 5, since i is shared
    }()
}
```
Fix (works in all versions):
```go
for i := 0; i < 5; i++ {
    i := i // shadow copy
    go func() {
        fmt.Println(i)
    }()
}
```
Or pass as parameter:
```go
go func(i int) {
    fmt.Println(i)
}(i)
```
> Note: Since **Go 1.22**, each loop iteration gets its own `i`, so the bug no longer occurs by default — but it's still important to know for older code and for interviews.

---

## 2. Channels

Channels let goroutines communicate and synchronize, without shared memory. "Do not communicate by sharing memory; share memory by communicating."

```go
ch := make(chan int)         // unbuffered channel
ch := make(chan int, 5)      // buffered channel, capacity 5

ch <- 10        // send
v := <-ch       // receive
v, ok := <-ch   // ok is false if channel closed and drained
close(ch)       // close a channel (only sender should close)
```

**Unbuffered vs Buffered:**
| Type | Behavior |
|---|---|
| Unbuffered | Send blocks until a receiver is ready (synchronous handoff) |
| Buffered | Send blocks only when buffer is full; receive blocks only when empty |

**Directional channels** (restrict usage, improve safety):
```go
func send(ch chan<- int) { ch <- 1 }   // send-only
func recv(ch <-chan int) { <-ch }      // receive-only
```

**Ranging over a channel:**
```go
for v := range ch {
    fmt.Println(v)
} // loop ends automatically when channel is closed
```

**Rules to remember:**
- Sending on a closed channel → **panic**.
- Closing an already-closed channel → **panic**.
- Closing a nil channel → **panic**.
- Receiving from a closed channel returns the zero value immediately (never blocks).
- Only the sender should close a channel, never the receiver.

---

## 3. select Statement

`select` lets a goroutine wait on multiple channel operations.

```go
select {
case msg1 := <-ch1:
    fmt.Println("received", msg1)
case msg2 := <-ch2:
    fmt.Println("received", msg2)
case ch3 <- 5:
    fmt.Println("sent to ch3")
default:
    fmt.Println("no channel ready") // non-blocking select
}
```

- If multiple cases are ready, one is chosen **pseudo-randomly**.
- `default` makes `select` non-blocking.
- Common pattern — timeout:
```go
select {
case res := <-resultChan:
    fmt.Println(res)
case <-time.After(2 * time.Second):
    fmt.Println("timeout")
}
```

---

## 4. sync Package

### 4.1 WaitGroup
Waits for a collection of goroutines to finish.
```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(i int) {
        defer wg.Done()
        fmt.Println(i)
    }(i)
}
wg.Wait() // blocks until counter hits 0
```
- `Add(n)` before starting goroutines (not inside them, to avoid race).
- `Done()` == `Add(-1)`.
- `Wait()` blocks until counter is 0.

### 4.2 Mutex / RWMutex
Protect shared state from concurrent access (data races).
```go
var mu sync.Mutex
var counter int

func increment() {
    mu.Lock()
    defer mu.Unlock()
    counter++
}
```
- `sync.RWMutex`: multiple readers OR one writer at a time.
  - `RLock()/RUnlock()` for reads, `Lock()/Unlock()` for writes.

### 4.3 Once
Runs a function exactly once, no matter how many goroutines call it (great for singleton init).
```go
var once sync.Once
once.Do(func() {
    fmt.Println("init")
})
```

### 4.4 sync.Map
A concurrency-safe map (useful for high read/write concurrency with disjoint key sets). Usually a plain `map` + `Mutex` is preferred unless a specific `sync.Map` use case fits (mostly-read, keys stable).

### 4.5 atomic package
For lock-free simple counters:
```go
var counter int64
atomic.AddInt64(&counter, 1)
atomic.LoadInt64(&counter)
```

---

## 5. Context Package

`context.Context` carries deadlines, cancellation signals, and request-scoped values across goroutines/API boundaries.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go worker(ctx)

func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("cancelled:", ctx.Err())
            return
        default:
            // do work
        }
    }
}
```

Variants:
- `context.WithCancel(parent)` — manual cancel
- `context.WithTimeout(parent, d)` — auto-cancel after duration
- `context.WithDeadline(parent, t)` — auto-cancel at a specific time
- `context.WithValue(parent, key, val)` — attach request-scoped data (avoid overusing for business logic)

---

## 6. Common Concurrency Patterns

### 6.1 Worker Pool

Producer puts jobs into a jobs channel

Workers read jobs from the channel

Each worker processes a job

Results go into a results channel

Main goroutine collects results
```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)

    for a := 1; a <= 5; a++ {
        fmt.Println(<-results)
    }
}
```

### 6.2 Fan-out / Fan-in
- **Fan-out**: multiple goroutines read from the same channel to parallelize work.
- **Fan-in**: multiple channels merged into a single channel.

```go
func fanIn(chs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(chs))
    for _, c := range chs {
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(c)
    }
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

### 6.3 Pipeline
Chain of stages connected by channels, each stage takes input channel and returns output channel.
```go
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}
```

### 6.4 Semaphore (limit concurrency using buffered channel)
```go
sem := make(chan struct{}, 3) // max 3 concurrent

for _, task := range tasks {
    sem <- struct{}{}
    go func(t Task) {
        defer func() { <-sem }()
        process(t)
    }(task)
}
```

---

## 7. Data Races & Detection

A **data race** happens when 2+ goroutines access the same variable concurrently and at least one is a write, without synchronization.

Detect with:
```bash
go run -race main.go
go test -race ./...
```

Avoid races by:
- Using channels to pass ownership of data.
- Using `Mutex`/`RWMutex` around shared state.
- Using `atomic` for simple counters.
- Never writing to a shared variable from multiple goroutines without sync.

---

## 8. Deadlocks & Common Pitfalls

- **Deadlock**: all goroutines are blocked waiting on each other. Go runtime detects *some* deadlocks (all goroutines asleep) and panics: `fatal error: all goroutines are asleep - deadlock!`
- **Goroutine leak**: a goroutine blocked forever (e.g., sending to a channel nobody reads) never gets garbage collected — leaks memory.
- **Forgetting to close channels** when using `range` over them → receiver blocks forever.
- **Closing a channel from the receiver side** → can cause "send on closed channel" panics elsewhere.
- **Sharing a `sync.WaitGroup` by value** (should always pass by pointer, or as struct field, never copy it).
- **Calling `wg.Add()` inside a goroutine** instead of before starting it → race between `Add` and `Wait`.

---

## 9. Practice Questions & Answers

### Q1. What is the difference between a goroutine and an OS thread?
**A:** A goroutine is a lightweight, user-space unit of concurrency managed by the Go runtime scheduler; it has a small growable stack (starts ~2KB) and thousands can run cheaply. An OS thread is managed by the operating system, has a larger fixed stack (often MBs), and context switches are more expensive. Go multiplexes many goroutines onto few OS threads (M:N scheduling).

---

### Q2. What happens if `main()` finishes before a spawned goroutine completes?
**A:** The program exits immediately; any goroutines still running are terminated without finishing. You must use synchronization (`sync.WaitGroup`, channels) to wait for goroutines to complete before `main` returns.

---

### Q3. What's the difference between buffered and unbuffered channels?
**A:** An unbuffered channel has zero capacity — a send blocks until a receiver is ready to receive (a synchronous "rendezvous"). A buffered channel has capacity `n`; a send only blocks once the buffer is full, and a receive only blocks when it's empty.

---

### Q4. What happens when you send to a closed channel? Receive from a closed channel?
**A:** Sending to a closed channel causes a **panic**. Receiving from a closed channel does **not** panic — it returns the channel's zero value immediately, and the second return value (`v, ok := <-ch`) will be `false` once the channel is closed and empty.

---

### Q5. Who should close a channel — the sender or receiver?
**A:** Only the **sender** should close a channel. Closing from the receiver side is a common bug source, since the sender may try to send afterward and panic ("send on closed channel").

---

### Q6. Write code to run 3 goroutines and wait for all to finish.
**A:**
```go
var wg sync.WaitGroup
for i := 1; i <= 3; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fmt.Println("worker", id)
    }(i)
}
wg.Wait()
```

---

### Q7. What is a data race, and how do you detect one in Go?
**A:** A data race occurs when two or more goroutines access the same memory location concurrently, and at least one access is a write, with no synchronization. Detect using Go's built-in race detector: `go run -race main.go` or `go test -race ./...`.

---

### Q8. Explain `sync.Mutex` vs `sync.RWMutex`.
**A:** `sync.Mutex` allows only one goroutine to hold the lock at a time (for both reads and writes). `sync.RWMutex` allows multiple concurrent readers (`RLock`/`RUnlock`) as long as no writer holds the lock, but only one writer (`Lock`/`Unlock`) at a time and writers block all readers. Use `RWMutex` when reads vastly outnumber writes.

---

### Q9. What does `select` do, and what happens if multiple cases are ready?
**A:** `select` blocks until one of its channel operations (send or receive) can proceed. If multiple cases are ready simultaneously, Go picks one **pseudo-randomly** (uniformly) to avoid starvation. If none are ready and there's a `default` case, that runs immediately (making the select non-blocking); without `default`, it blocks.

---

### Q10. How do you implement a timeout for a channel operation?
**A:**
```go
select {
case res := <-ch:
    fmt.Println("got", res)
case <-time.After(3 * time.Second):
    fmt.Println("timed out")
}
```

---

### Q11. What is a goroutine leak? Give an example.
**A:** A goroutine leak is a goroutine that blocks forever and is never cleaned up, wasting memory. Example: a goroutine sends to an unbuffered channel that nobody ever receives from:
```go
func leak() {
    ch := make(chan int)
    go func() {
        ch <- 1 // blocks forever if nobody reads
    }()
    // function returns without reading from ch — goroutine leaked
}
```
Fix: use buffered channels, `context` cancellation, or ensure a receiver always exists.

---

### Q12. What's wrong with this code?
```go
func main() {
    for i := 0; i < 3; i++ {
        go func() {
            fmt.Println(i)
        }()
    }
    time.Sleep(time.Second)
}
```
**A:** (Go versions before 1.22) All goroutines share the same loop variable `i`, so by the time they run, `i` may already be `3`, printing `3 3 3` instead of `0 1 2`. Fix by shadowing (`i := i`) or passing `i` as a parameter. Since Go 1.22, each iteration has its own `i`, so this specific bug is fixed by default — but it's still an important gotcha to know for older codebases and interviews.

---

### Q13. How does `context.WithTimeout` differ from `context.WithCancel`?
**A:** `WithCancel` returns a context that is cancelled only when you explicitly call the returned `cancel()` function. `WithTimeout` (built on `WithDeadline`) automatically cancels the context after the specified duration elapses, in addition to allowing manual cancellation — you should still call `cancel()` (usually via `defer`) to release resources immediately once done.

---

### Q14. What does `runtime.GOMAXPROCS` control?
**A:** It sets the maximum number of OS threads that can execute Go code simultaneously (i.e., parallelism). By default it equals the number of logical CPUs. It does not limit the total number of goroutines you can create — only how many can run truly in parallel at once.

---

### Q15. How would you limit concurrency to at most N goroutines running at once?
**A:** Use a buffered channel as a semaphore:
```go
sem := make(chan struct{}, N)
for _, job := range jobs {
    sem <- struct{}{}          // acquire
    go func(j Job) {
        defer func() { <-sem }() // release
        process(j)
    }(job)
}
```

---

### Q16. What is the output (or issue) with this code?
```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    fmt.Println("hello")
}()
wg.Wait()
```
**A:** This deadlocks (or `Wait()` blocks forever), because `wg.Done()` is never called inside the goroutine. `Add(1)` increments the counter but nothing decrements it. Fix: add `defer wg.Done()` inside the goroutine.

---

### Q17. Difference between `sync.Once` use case and simply checking a boolean flag?
**A:** A boolean flag check-and-set (`if !initialized { initialized = true; init() }`) is **not** goroutine-safe — multiple goroutines could pass the check simultaneously before the flag is set, causing `init()` to run multiple times. `sync.Once.Do()` guarantees the function runs exactly once across all goroutines, safely, using internal synchronization.

---

### Q18. How do you merge (fan-in) multiple channels into one?
**A:**
```go
func fanIn(chs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(chs))
    for _, c := range chs {
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(c)
    }
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

---

### Q19. What's the difference between concurrency and parallelism?
**A:** Concurrency is about **structuring** a program to handle multiple tasks that can make progress independently (they may or may not run at the exact same instant). Parallelism is about **actually executing** multiple tasks at the same literal moment, typically requiring multiple CPU cores. Go's goroutines give you concurrency; whether they run in parallel depends on `GOMAXPROCS` and the number of available cores.

---

### Q20. What happens if you copy a `sync.Mutex` or `sync.WaitGroup` (e.g., pass by value into a function)?
**A:** This is a bug. Both types contain internal state that must not be copied after first use — copying creates a separate lock/counter that no longer synchronizes with the original, breaking mutual exclusion or wait semantics. `go vet` flags this. Always pass mutexes/waitgroups by pointer, or embed them in a struct that is itself passed by pointer.

---

## 10. Quick Reference Cheat-Sheet

| Concept | Keyword/Type | Notes |
|---|---|---|
| Start goroutine | `go f()` | Non-blocking, no return value access |
| Wait for goroutines | `sync.WaitGroup` | `Add`, `Done`, `Wait` |
| Mutual exclusion | `sync.Mutex` / `RWMutex` | Lock/Unlock, RLock/RUnlock |
| Run once | `sync.Once` | `.Do(func())` |
| Channel | `chan T` | `make(chan T)`, `make(chan T, n)` for buffered |
| Multiplex channels | `select` | Random pick if multiple ready; `default` = non-blocking |
| Cancellation | `context.Context` | `WithCancel`, `WithTimeout`, `WithDeadline`, `WithValue` |
| Atomic ops | `sync/atomic` | Lock-free counters |
| Race detection | `go run -race` | Always test concurrent code with `-race` |

---

*End of notes. Practice by writing worker pools, pipelines, and rate limiters from scratch without looking at the examples above — that's the fastest way to internalize these patterns.*

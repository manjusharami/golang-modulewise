# Go Concurrency — Interview Prep Notes

A deep, interview-focused reference on Go concurrency: concepts, gotchas, code patterns, and the questions interviewers actually ask.

---

## 1. Foundations you must be able to explain clearly

### Goroutines
- A goroutine is a function running concurrently, managed by the **Go runtime scheduler**, not the OS.
- Started with `go f()`. Starts with a tiny stack (~2KB) that grows/shrinks dynamically — this is why you can spawn thousands cheaply, unlike OS threads (~1–8MB fixed stacks).
- The Go scheduler uses an **M:N model**: M goroutines multiplexed onto N OS threads (`GOMAXPROCS` controls how many threads run Go code simultaneously).
- **Interview trap**: "Does `go f()` guarantee `f` runs before the next line?" — No. It's scheduled, not immediate. Order is not guaranteed.
- **Interview trap**: if `main()` returns, all goroutines are killed immediately — no cleanup, no waiting.

### Channels
- Typed conduits for communication: `ch := make(chan int)` (unbuffered) or `make(chan int, n)` (buffered).
- Unbuffered channel = **synchronous handoff**: send blocks until a receiver is ready, and vice versa.
- Buffered channel: send blocks only when the buffer is full; receive blocks only when it's empty.
- Directional types: `chan<- int` (send-only), `<-chan int` (receive-only) — used in function signatures to express intent and catch misuse at compile time.
- `close(ch)` signals "no more values." Reading a closed channel yields the zero value immediately with `ok == false`.

**Golden rules interviewers probe:**
1. Only the **sender** should close a channel, never the receiver.
2. Closing an already-closed channel **panics**.
3. Sending on a closed channel **panics**.
4. Receiving from a closed channel never blocks (returns zero value, `ok=false`).
5. A `nil` channel blocks forever on send and receive (used deliberately in `select` to "disable" a case).

### select
- Lets a goroutine wait on multiple channel operations at once; runs whichever case is ready.
- If multiple cases are ready, Go picks one **pseudo-randomly** (not first-come order) — a common quiz question.
- `default` case makes the whole `select` non-blocking.
- Classic timeout pattern:
```go
select {
case res := <-resultCh:
    fmt.Println(res)
case <-time.After(2 * time.Second):
    fmt.Println("timed out")
}
```

### sync package
| Tool | Purpose |
|---|---|
| `sync.WaitGroup` | Wait for a batch of goroutines to finish (`Add`/`Done`/`Wait`) |
| `sync.Mutex` | Mutual exclusion lock for shared state |
| `sync.RWMutex` | Multiple readers OR one writer |
| `sync.Once` | Run initialization code exactly once, safely, across goroutines |
| `sync.Map` | Concurrent-safe map for specific high-read/low-write patterns (rarely the default choice — interviewers like asking *why not just use a map + mutex*) |
| `sync/atomic` | Lock-free primitives for simple counters/flags |

**Rule of thumb interviewers want to hear**: *"Don't communicate by sharing memory; share memory by communicating."* Use channels to pass ownership of data between goroutines; use mutexes to protect state that's genuinely shared and accessed concurrently. Neither is universally "better" — know when each fits.

### context
- `context.Background()` / `context.TODO()` are roots.
- `context.WithCancel(ctx)`, `WithTimeout(ctx, d)`, `WithDeadline(ctx, t)` derive a cancellable child.
- Convention: `ctx` is always the **first parameter**, never stored in a struct.
- Cancellation propagates down the tree: cancelling a parent cancels all children.
- `ctx.Done()` returns a channel closed when the context is cancelled/expired — check it in `select` to abort work early.
- `ctx.Err()` tells you why it's done (`context.Canceled` or `context.DeadlineExceeded`).
- `ctx.Value()` exists for request-scoped metadata (e.g. trace IDs) — interviewers expect you to say it should **not** be used to pass optional function parameters.

---

## 2. Frequently asked interview questions (with answers)

**Q: What's the difference between concurrency and parallelism?**
Concurrency is about structuring a program to deal with many things at once (tasks can overlap in time, e.g. via interleaving on a single core). Parallelism is actually doing many things at the same time (needs multiple cores). Go's runtime makes concurrent code easy to write; whether it runs in parallel depends on `GOMAXPROCS` and available cores.

**Q: How does the Go scheduler work at a high level?**
It's a cooperative, work-stealing scheduler (G-M-P model): **G**oroutines, **M**achine threads, **P**rocessors (logical, count = `GOMAXPROCS`). Each P has a local run queue of goroutines; idle Ps steal work from busy ones. Goroutines yield at function calls, channel ops, syscalls, and GC safepoints — not truly preemptible everywhere before Go 1.14, but async preemption was added in 1.14+ to handle tight loops.

**Q: What causes a deadlock, and how do you detect one?**
All goroutines are blocked waiting on each other with nothing left to make progress (e.g. two goroutines each waiting to receive on a channel the other should send on, or a mutex locked twice by the same goroutine). The Go runtime detects when **all** goroutines are asleep and panics with `fatal error: all goroutines are asleep - deadlock!` at compile-run time. Partial deadlocks (some but not all goroutines stuck) are *not* automatically detected — that's a goroutine leak, and you need profiling (`pprof`) or the `-race` detector's cousin tools to catch it.

**Q: What's a goroutine leak, and give an example.**
A goroutine that's blocked forever and never returns, silently leaking memory. Classic example: sending on an unbuffered channel that nobody is left to receive from, because the receiver already returned:
```go
func leaky() {
    ch := make(chan int)
    go func() {
        ch <- compute() // blocks forever if nobody reads ch
    }()
    // function returns without ever reading from ch
}
```
Fix: use a buffered channel sized to avoid blocking, add a `context` for cancellation, or make sure every send has a guaranteed receiver.

**Q: How do you detect race conditions?**
Run tests/binary with `go run -race` or `go test -race`. It instruments memory accesses at compile time and flags unsynchronized concurrent read/write to the same location. Always mention: the race detector only catches races that actually occur during that run — it's not a static guarantee for all possible interleavings.

**Q: Explain the classic loop-variable-capture bug (very common interview question).**
Before Go 1.22, this printed unexpected values:
```go
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i) // captures the loop variable by reference — often prints 3,3,3
    }()
}
```
Because all closures captured the *same* `i`. Fixes pre-1.22:
```go
for i := 0; i < 3; i++ {
    i := i // shadow: new variable per iteration
    go func() { fmt.Println(i) }()
}
// or pass as an argument:
go func(i int) { fmt.Println(i) }(i)
```
**Go 1.22+ changed loop semantics** so each iteration gets its own `i` automatically — know this, and mention which Go version you're assuming.

**Q: When would you use a buffered vs unbuffered channel?**
Unbuffered when you want a strict handshake / backpressure — the sender knows the receiver has "accepted" the value. Buffered when you want to decouple producer and consumer speed up to some limit (e.g. a queue of pending jobs), or to avoid blocking a fast producer for short bursts. Buffered channels don't remove the *need* for correct shutdown logic — they just add slack.

**Q: How would you implement a worker pool?**
```go
func workerPool(jobs <-chan int, results chan<- int, numWorkers int) {
    var wg sync.WaitGroup
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs { // exits when jobs channel is closed & drained
                results <- process(job)
            }
        }()
    }
    go func() {
        wg.Wait()
        close(results)
    }()
}
```
Key interview points: workers pull from a shared jobs channel (`range` exits automatically on close), a separate goroutine waits on the `WaitGroup` before closing `results` to avoid closing it while workers still write to it.

**Q: How do you implement fan-out / fan-in?**
Fan-out: multiple goroutines read from the same input channel. Fan-in: multiple goroutines write to the same output channel, merged into one stream.
```go
func fanIn(cs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(cs))
    for _, c := range cs {
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

**Q: How do you cancel a long-running goroutine cleanly?**
Pass a `context.Context`, check `ctx.Done()` in the goroutine's loop/select, and return promptly when it fires. Never use a raw `bool` flag checked without synchronization — that's itself a race.
```go
func worker(ctx context.Context, jobs <-chan int) {
    for {
        select {
        case <-ctx.Done():
            return
        case j, ok := <-jobs:
            if !ok {
                return
            }
            process(j)
        }
    }
}
```

**Q: What's the difference between `sync.Mutex` and `sync.RWMutex`? When would RWMutex hurt you?**
`RWMutex` allows many concurrent readers or one writer. It helps when reads vastly outnumber writes. It can hurt when: critical sections are very short (the extra bookkeeping overhead can exceed the benefit of a plain `Mutex`), or when writers get starved under heavy read load (implementation-dependent).

**Q: What does `sync.Once` guarantee, and what's a real use case?**
Exactly-once execution of a function, safe across goroutines, typically for lazy singleton initialization (e.g. a global config or connection pool):
```go
var once sync.Once
var instance *Config
func GetConfig() *Config {
    once.Do(func() { instance = loadConfig() })
    return instance
}
```

**Q: Why might you choose `sync/atomic` over a mutex?**
For very simple operations (increment a counter, swap a pointer, flip a flag) atomic operations avoid lock overhead and are lock-free. They don't generalize to protecting multi-field structs or compound invariants — that's when you need a mutex instead.

**Q: What happens if you range over a channel that's never closed?**
The `range` loop blocks forever waiting for the next value once nothing more is sent — a leak. Always ensure the sender closes the channel when done, and only the sender closes it.

**Q: How would you rate-limit work across goroutines?**
A common pattern is a **buffered channel as a semaphore**, or `time.Ticker` / `golang.org/x/time/rate` for actual rate limiting:
```go
sem := make(chan struct{}, maxConcurrent)
for _, job := range jobs {
    sem <- struct{}{} // acquire
    go func(j Job) {
        defer func() { <-sem }() // release
        process(j)
    }(job)
}
```

**Q: What's the difference between `context.WithCancel` and `context.WithTimeout`?**
`WithCancel` gives you an explicit `cancel()` function you call yourself. `WithTimeout` (built on `WithCancel` internally) auto-cancels after a duration. Both must have their `cancel` function called (usually via `defer cancel()`) to release resources promptly, even if the operation finishes early — otherwise the context and its timer leak until the deadline fires.

**Q: Can two goroutines write to the same map concurrently?**
No — Go's built-in `map` is not safe for concurrent writes (or a concurrent write + read). Doing so without synchronization panics at runtime with `fatal error: concurrent map writes` (a deliberate, checked failure, not silent corruption... though reads-during-write can still corrupt). Fix with a `sync.Mutex`/`RWMutex` wrapper, or `sync.Map` for specific access patterns (many reads, keys written once, or disjoint key sets across goroutines).

---

## 3. Gotchas interviewers love to test

- **Unbuffered channel + no receiver** → deadlock.
- **Closing a channel with active senders** → panic on next send.
- **Double close** → panic.
- **Nil channel in select** → that case is effectively disabled (useful trick, not a bug, when intentional).
- **Forgetting `defer wg.Done()`** → `Wait()` blocks forever if a goroutine panics or returns early before calling `Done()`.
- **Passing a WaitGroup by value instead of pointer** → each goroutine gets its own copy, `Wait()` never sees the right count.
- **Mutex copied by value** (e.g. embedding a `sync.Mutex` in a struct passed by value) → each copy has independent lock state, defeats the purpose. `go vet` catches this.
- **Data race on a slice/map without a mutex** → silent corruption or panic, use `-race` to catch it in testing.

---

## 4. How to structure your answer in an interview

1. **State the correct behavior first** (what the code does / should do).
2. **Name the underlying mechanism** (scheduler, channel semantics, memory model) — this signals depth beyond memorized syntax.
3. **Give a minimal code example** if asked "how would you implement...".
4. **Mention the failure mode** (deadlock, leak, race, panic) and how you'd detect it (`-race`, `pprof`, `go vet`).
5. If genuinely unsure of an edge case (e.g. exact scheduler internals), say what you know and be explicit about the boundary of your certainty — interviewers respect that far more than a confident guess.

---

## 5. Quick-fire review (say the answer out loud before checking)

- What does `close()` do to pending receivers on a channel? *(unblocks them immediately with zero value, ok=false)*
- What's the zero value of a channel, and what happens if you use it? *(nil; sends/receives block forever)*
- Does `go func(){}()` guarantee the goroutine starts immediately? *(no)*
- Name two ways to avoid the pre-1.22 loop variable capture bug. *(shadow the variable inside the loop; pass it as a function argument)*
- What error does `ctx.Err()` return after a timeout fires? *(`context.DeadlineExceeded`)*
- Is a Go map safe for concurrent read/write without synchronization? *(no)*
- What tool flags data races at test/run time? *(`-race`)*

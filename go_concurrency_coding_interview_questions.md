# Go Concurrency — Coding Interview Questions (with Solutions)

Hands-on coding problems, ordered roughly easy → hard. Try to solve each one yourself first, then check the solution.

---

## Q1. Print numbers 1 to 10 using two goroutines alternately (one prints odd, one prints even)

**Problem:** Use channels to make two goroutines take turns printing, so output is strictly `1,2,3,4,...,10` in order.

**Solution:**
```go
package main

import "fmt"

func main() {
	oddCh := make(chan struct{})
	evenCh := make(chan struct{})
	done := make(chan struct{})

	go func() { // odd printer
		for i := 1; i <= 10; i += 2 {
			<-oddCh
			fmt.Println(i)
			if i+1 <= 10 {
				evenCh <- struct{}{}
			} else {
				done <- struct{}{}
			}
		}
	}()

	go func() { // even printer
		for i := 2; i <= 10; i += 2 {
			<-evenCh
			fmt.Println(i)
			if i+1 <= 10 {
				oddCh <- struct{}{}
			} else {
				done <- struct{}{}
			}
		}
	}()

	oddCh <- struct{}{} // kick off
	<-done
}
```
**Key idea:** unbuffered channels act as a "baton" passed between goroutines to enforce ordering.

---

## Q2. Implement a thread-safe counter

**Problem:** Multiple goroutines increment a shared counter 1000 times each. Ensure the final value is correct.

**Solution (Mutex):**
```go
type Counter struct {
	mu    sync.Mutex
	count int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func main() {
	c := &Counter{}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				c.Inc()
			}
		}()
	}
	wg.Wait()
	fmt.Println(c.count) // 10000
}
```
**Alternative (lock-free, faster):**
```go
var count int64
atomic.AddInt64(&count, 1)
```

---

## Q3. Worker Pool: process N jobs with a fixed number of workers

**Problem:** Given 20 jobs and 4 workers, distribute jobs and collect results.

**Solution:**
```go
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		results <- j * j
	}
}

func main() {
	const numJobs, numWorkers = 20, 4
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(results)

	for r := range results {
		fmt.Println(r)
	}
}
```

---

## Q4. Rate Limiter — allow at most 1 request per 200ms

**Problem:** Implement a simple rate limiter using `time.Ticker`.

**Solution:**
```go
func main() {
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	limiter := time.NewTicker(200 * time.Millisecond)
	defer limiter.Stop()

	for req := range requests {
		<-limiter.C // block until next tick
		fmt.Println("processing request", req, time.Now())
	}
}
```
**Burst variant:** use a buffered channel refilled by a ticker to allow short bursts.

---

## Q5. Implement a Semaphore to limit concurrent goroutines to N

**Problem:** Run 10 tasks but never more than 3 concurrently.

**Solution:**
```go
func main() {
	sem := make(chan struct{}, 3)
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem <- struct{}{}        // acquire
			defer func() { <-sem }() // release

			fmt.Println("task", id, "running")
			time.Sleep(500 * time.Millisecond)
		}(i)
	}
	wg.Wait()
}
```

---

## Q6. Detect and fix a deadlock

**Problem:** Find the bug.
```go
func main() {
	ch := make(chan int)
	ch <- 1        // blocks forever, no receiver
	fmt.Println(<-ch)
}
```
**Why it deadlocks:** `ch` is unbuffered; the send blocks waiting for a receiver, but the receive line never executes because we're stuck on the line above (single goroutine — main itself).

**Fix options:**
```go
// Option A: buffered channel
ch := make(chan int, 1)
ch <- 1
fmt.Println(<-ch)

// Option B: separate goroutine for send
ch := make(chan int)
go func() { ch <- 1 }()
fmt.Println(<-ch)
```

---

## Q7. Fan-in: merge multiple channels into one

**Problem:** Given `chan int` producers `c1, c2, c3`, merge all values into a single output channel.

**Solution:**
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

---

## Q8. Pipeline: generate → square → sum

**Problem:** Build a 3-stage pipeline using channels.

**Solution:**
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

func sum(in <-chan int) int {
	total := 0
	for n := range in {
		total += n
	}
	return total
}

func main() {
	total := sum(square(gen(1, 2, 3, 4, 5)))
	fmt.Println(total) // 55
}
```

---

## Q9. Timeout a long-running operation

**Problem:** Call a function that may take arbitrarily long; give up after 2 seconds.

**Solution:**
```go
func slowOp(ch chan<- string) {
	time.Sleep(3 * time.Second)
	ch <- "done"
}

func main() {
	ch := make(chan string, 1)
	go slowOp(ch)

	select {
	case res := <-ch:
		fmt.Println(res)
	case <-time.After(2 * time.Second):
		fmt.Println("timeout: operation took too long")
	}
}
```
**With context (preferred in real APIs):**
```go
func slowOpCtx(ctx context.Context, ch chan<- string) {
	select {
	case <-time.After(3 * time.Second):
		ch <- "done"
	case <-ctx.Done():
		return
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch := make(chan string, 1)
	go slowOpCtx(ctx, ch)

	select {
	case res := <-ch:
		fmt.Println(res)
	case <-ctx.Done():
		fmt.Println("timeout:", ctx.Err())
	}
}
```

---

## Q10. Producer-Consumer with bounded buffer

**Problem:** One producer generates numbers, multiple consumers process them, buffer size 5.

**Solution:**
```go
func producer(ch chan<- int, n int) {
	for i := 1; i <= n; i++ {
		ch <- i
		fmt.Println("produced", i)
	}
	close(ch)
}

func consumer(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {
		fmt.Printf("consumer %d got %d\n", id, v)
	}
}

func main() {
	ch := make(chan int, 5) // bounded buffer
	var wg sync.WaitGroup

	go producer(ch, 20)

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go consumer(i, ch, &wg)
	}
	wg.Wait()
}
```

---

## Q11. Dining Philosophers (classic deadlock-avoidance problem)

**Problem:** 5 philosophers, 5 forks, avoid deadlock and starvation.

**Solution (asymmetric ordering to break circular wait):**
```go
type Fork struct{ sync.Mutex }

func philosopher(id int, left, right *Fork, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 3; i++ { // eat 3 times
		// Break symmetry: last philosopher picks up right fork first
		first, second := left, right
		if id == 4 {
			first, second = right, left
		}
		first.Lock()
		second.Lock()

		fmt.Printf("philosopher %d eating\n", id)
		time.Sleep(100 * time.Millisecond)

		second.Unlock()
		first.Unlock()
	}
}

func main() {
	forks := make([]*Fork, 5)
	for i := range forks {
		forks[i] = &Fork{}
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go philosopher(i, forks[i], forks[(i+1)%5], &wg)
	}
	wg.Wait()
}
```
**Key idea:** deadlock happens when everyone grabs their left fork simultaneously, then waits forever for the right. Breaking the symmetry (one philosopher picks up in reverse order) prevents the circular wait condition.

---

## Q12. Implement `sync.WaitGroup`-like behavior using only channels

**Problem:** Don't use `sync.WaitGroup`; wait for N goroutines using channels only.

**Solution:**
```go
func main() {
	n := 5
	done := make(chan struct{}, n)

	for i := 0; i < n; i++ {
		go func(id int) {
			fmt.Println("working", id)
			done <- struct{}{}
		}(i)
	}

	for i := 0; i < n; i++ {
		<-done
	}
	fmt.Println("all done")
}
```

---

## Q13. Print "ping" and "pong" alternately N times using channels

**Solution:**
```go
func main() {
	ping := make(chan struct{})
	pong := make(chan struct{})
	n := 5

	go func() {
		for i := 0; i < n; i++ {
			<-ping
			fmt.Println("ping")
			pong <- struct{}{}
		}
	}()

	go func() {
		for i := 0; i < n; i++ {
			<-pong
			fmt.Println("pong")
			if i < n-1 {
				ping <- struct{}{}
			}
		}
	}()

	ping <- struct{}{}
	time.Sleep(time.Second) // or use a done channel in real code
}
```

---

## Q14. Find first error among N concurrent operations (fail fast)

**Problem:** Run N tasks concurrently; return as soon as any one fails, without waiting for the rest.

**Solution:**
```go
func doWork(id int) error {
	time.Sleep(time.Duration(id) * 100 * time.Millisecond)
	if id == 3 {
		return fmt.Errorf("task %d failed", id)
	}
	return nil
}

func main() {
	errCh := make(chan error, 1)
	n := 5

	for i := 1; i <= n; i++ {
		go func(id int) {
			if err := doWork(id); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(i)
	}

	select {
	case err := <-errCh:
		fmt.Println("failed fast:", err)
	case <-time.After(2 * time.Second):
		fmt.Println("all succeeded (or timeout)")
	}
}
```
**Better version with `errgroup` (golang.org/x/sync/errgroup):**
```go
g, ctx := errgroup.WithContext(context.Background())
for i := 1; i <= n; i++ {
	i := i
	g.Go(func() error {
		return doWork(ctx, i)
	})
}
if err := g.Wait(); err != nil {
	fmt.Println("first error:", err)
}
```

---

## Q15. Implement a debounce function for concurrent calls

**Problem:** Given rapid repeated calls, only execute the action after calls stop coming for 300ms.

**Solution:**
```go
func debounce(d time.Duration, f func()) func() {
	var mu sync.Mutex
	var timer *time.Timer

	return func() {
		mu.Lock()
		defer mu.Unlock()
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(d, f)
	}
}

func main() {
	debounced := debounce(300*time.Millisecond, func() {
		fmt.Println("action executed")
	})

	for i := 0; i < 5; i++ {
		debounced()
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(time.Second)
}
```

---

## Q16. Print numbers 1..100, but goroutine A prints multiples of 3, goroutine B prints the rest (in order)

**Problem:** Tests ordering + coordination, common "FizzBuzz concurrency" variant.

**Solution:**
```go
func main() {
	turnA := make(chan struct{}, 1)
	turnB := make(chan struct{}, 1)
	done := make(chan struct{})
	n := 100

	go func() { // multiples of 3
		for i := 1; i <= n; i++ {
			if i%3 == 0 {
				<-turnA
				fmt.Println("A:", i)
				turnB <- struct{}{}
			}
		}
	}()

	go func() { // everything else
		for i := 1; i <= n; i++ {
			if i%3 != 0 {
				<-turnB
				fmt.Println("B:", i)
				if i < n {
					turnA <- struct{}{}
				} else {
					done <- struct{}{}
				}
			}
		}
	}()

	turnB <- struct{}{} // 1 is not a multiple of 3, B starts
	<-done
}
```
*(Note: real interviews often accept a simpler version that doesn't need strict global ordering — clarify requirements before over-engineering.)*

---

## Q17. Safe map with concurrent reads/writes

**Problem:** Implement a goroutine-safe map without using `sync.Map`.

**Solution:**
```go
type SafeMap struct {
	mu   sync.RWMutex
	data map[string]int
}

func NewSafeMap() *SafeMap {
	return &SafeMap{data: make(map[string]int)}
}

func (m *SafeMap) Set(key string, val int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
}

func (m *SafeMap) Get(key string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}
```

---

## Q18. Cancel a group of goroutines using context

**Problem:** Start 5 workers; cancel all of them after 1 second.

**Solution:**
```go
func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d stopping: %v\n", id, ctx.Err())
			return
		default:
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("worker %d working\n", id)
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(ctx, id)
		}(i)
	}
	wg.Wait()
}
```

---

## Q19. Detect a data race — spot the bug

**Problem:** What's wrong here, and how do you fix it?
```go
func main() {
	m := make(map[int]int)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m[i] = i * i // concurrent map write — DATA RACE / possible panic
		}(i)
	}
	wg.Wait()
}
```
**Why it's broken:** plain Go maps are **not** safe for concurrent writes; this can panic with `fatal error: concurrent map writes` or corrupt data. Run with `go run -race` to confirm.

**Fix (Mutex):**
```go
var mu sync.Mutex
m := make(map[int]int)
...
go func(i int) {
	defer wg.Done()
	mu.Lock()
	m[i] = i * i
	mu.Unlock()
}(i)
```
**Fix (sync.Map):**
```go
var m sync.Map
...
m.Store(i, i*i)
```

---

## Q20. Implement a simple concurrent-safe LRU cache (bonus, harder)

**Problem:** Basic thread-safe LRU with `Get`/`Put`, fixed capacity.

**Solution:**
```go
type node struct {
	key, val   int
	prev, next *node
}

type LRUCache struct {
	mu       sync.Mutex
	cap      int
	m        map[int]*node
	head, tail *node // head = most recent, tail = least recent
}

func NewLRUCache(capacity int) *LRUCache {
	head, tail := &node{}, &node{}
	head.next, tail.prev = tail, head
	return &LRUCache{cap: capacity, m: make(map[int]*node), head: head, tail: tail}
}

func (c *LRUCache) remove(n *node) {
	n.prev.next, n.next.prev = n.next, n.prev
}

func (c *LRUCache) insertFront(n *node) {
	n.next = c.head.next
	n.prev = c.head
	c.head.next.prev = n
	c.head.next = n
}

func (c *LRUCache) Get(key int) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	n, ok := c.m[key]
	if !ok {
		return 0, false
	}
	c.remove(n)
	c.insertFront(n)
	return n.val, true
}

func (c *LRUCache) Put(key, val int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if n, ok := c.m[key]; ok {
		n.val = val
		c.remove(n)
		c.insertFront(n)
		return
	}
	if len(c.m) >= c.cap {
		lru := c.tail.prev
		c.remove(lru)
		delete(c.m, lru.key)
	}
	n := &node{key: key, val: val}
	c.m[key] = n
	c.insertFront(n)
}
```

---

## Quick Tips for the Interview

- Always mention `go run -race` when discussing correctness of concurrent code.
- Explain **why** unbuffered channels synchronize (rendezvous point) vs buffered channels (async up to capacity).
- If asked "how would you test this," mention `-race`, table-driven tests, and stress-testing with high goroutine counts.
- Know the difference between goroutine leak (blocked forever) and deadlock (whole program stuck) — interviewers often ask you to spot one.
- If given an open-ended concurrency design question, always clarify: ordering requirements, expected concurrency level, and failure/cancellation behavior before coding.

---

*Practice tip: re-implement each solution from a blank file without peeking — that's what actually sticks for interviews.*

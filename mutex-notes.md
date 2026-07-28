# sync.Mutex — Detailed Notes & Practical Use

## What is it?

`Mutex` (mutual exclusion lock, from the `sync` package) ensures that **only one goroutine at a time** can access a critical section of code / shared data — preventing race conditions.

---

## How it works

| Method | Purpose |
|---|---|
| `mu.Lock()` | Acquires the lock. Blocks if another goroutine already holds it. |
| `mu.Unlock()` | Releases the lock, allowing another waiting goroutine to acquire it. |

### Key Rules

- Always `Unlock()` — typically with `defer` right after `Lock()` — to avoid deadlocks if a panic occurs.
- Never copy a `Mutex` after first use — pass by pointer or embed in a struct.
- Zero value is ready to use: `var mu sync.Mutex`
- Locking is **not reentrant** — a goroutine that already holds the lock will deadlock if it calls `Lock()` again.
- Keep the critical section as small as possible — don't do slow I/O while holding a lock.

### Zero value is usable
```go
var mu sync.Mutex // ready to use
```

---

## Live Example — Race Condition Without Mutex

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // ❌ race condition: read-modify-write not atomic
		}()
	}

	wg.Wait()
	fmt.Println("Final counter:", counter) // often NOT 1000
}
```

Run with `go run -race main.go` and it will report a **data race**.

---

## Fixed With Mutex

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			counter++ // ✅ safe now, only one goroutine at a time
		}()
	}

	wg.Wait()
	fmt.Println("Final counter:", counter) // always 1000
}
```

---

## RWMutex — Read/Write Lock Variant

Use when you have **many readers, few writers** — allows concurrent reads but exclusive writes.

| Method | Purpose |
|---|---|
| `mu.RLock()` / `mu.RUnlock()` | Multiple goroutines can hold read lock simultaneously |
| `mu.Lock()` / `mu.Unlock()` | Exclusive — blocks all readers and writers |

```go
type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func (c *Cache) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}
```

**Why it matters:** Reads (`Get`) don't block each other — only a `Set` blocks everything. Great for read-heavy caches/config stores.

---

## Where Mutex is Used in Production

### 1. Thread-safe in-memory cache

Real scenario: an API service caches computed results (e.g., pricing, session data) in a map shared across many request-handling goroutines.

```go
type PriceCache struct {
	mu    sync.RWMutex
	prices map[string]float64
}

func NewPriceCache() *PriceCache {
	return &PriceCache{prices: make(map[string]float64)}
}

func (c *PriceCache) Get(sku string) (float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	price, ok := c.prices[sku]
	return price, ok
}

func (c *PriceCache) Set(sku string, price float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prices[sku] = price
}
```

**Why it matters:** Go maps are NOT safe for concurrent read/write — this pattern is everywhere in production Go services (config caches, session stores, feature flags).

---

### 2. Protecting shared counters / metrics

Real scenario: tracking request counts, error counts, or rate-limiting state across concurrent HTTP handlers.

```go
type Metrics struct {
	mu       sync.Mutex
	requests int
	errors   int
}

func (m *Metrics) RecordRequest() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requests++
}

func (m *Metrics) RecordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors++
}
```

**Why it matters:** Used in middleware for request counters, error rates, and custom Prometheus-style metrics collectors.

---

### 3. Connection pool / resource manager

Real scenario: managing a limited pool of DB connections or worker slots shared across goroutines.

```go
type Pool struct {
	mu    sync.Mutex
	conns []*Connection
}

func (p *Pool) Get() *Connection {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.conns) == 0 {
		return nil // or create new
	}
	conn := p.conns[len(p.conns)-1]
	p.conns = p.conns[:len(p.conns)-1]
	return conn
}

func (p *Pool) Put(conn *Connection) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.conns = append(p.conns, conn)
}
```

**Why it matters:** Similar logic underlies real connection pools like `database/sql`'s internal pool.

---

### 4. Config hot-reload (RWMutex)

Real scenario: application config is reloaded periodically from a file/remote source, while many goroutines read it concurrently.

```go
type ConfigStore struct {
	mu  sync.RWMutex
	cfg Config
}

func (c *ConfigStore) Get() Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cfg
}

func (c *ConfigStore) Reload(newCfg Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cfg = newCfg
}
```

**Why it matters:** Lets thousands of request-handling goroutines read config with minimal contention, while updates happen safely in the background.

---

## Common Mistakes

```go
// ❌ Forgetting to unlock (deadlock risk)
mu.Lock()
if err != nil {
    return err // Unlock never called!
}
mu.Unlock()

// ✅ Fix: always defer immediately after Lock()
mu.Lock()
defer mu.Unlock()
if err != nil {
    return err
}
```

```go
// ❌ Copying a struct containing a Mutex
type Counter struct {
	mu    sync.Mutex
	count int
}

func increment(c Counter) { // BAD: copies the mutex!
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

// ✅ Fix: use a pointer receiver
func increment(c *Counter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}
```

---

## Mutex vs RWMutex vs Channels

| Use Case | Tool |
|---|---|
| Simple exclusive access to shared data | `Mutex` |
| Many reads, few writes | `RWMutex` |
| Passing ownership of data / signaling between goroutines | Channels |
| Protecting a counter/map/struct field | `Mutex` |
| Coordinating a pipeline or producer/consumer flow | Channels |

> Go proverb: *"Don't communicate by sharing memory; share memory by communicating."*
> In practice, both are used — Mutex for simple shared state, channels for coordination/data flow.

---

## Summary Cheat Sheet

```go
var mu sync.Mutex

mu.Lock()
defer mu.Unlock()
// critical section — access shared data here
```

```go
var mu sync.RWMutex

mu.RLock()
defer mu.RUnlock()
// read-only access

mu.Lock()
defer mu.Unlock()
// write access
```

**Production use cases recap:**
- Thread-safe in-memory caches (maps aren't safe by default)
- Shared counters/metrics
- Connection/resource pools
- Config hot-reload with concurrent reads (RWMutex)

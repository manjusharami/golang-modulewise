 Goroutine vs OS Thread

Managed by: Goroutines are managed by the Go runtime's user-space scheduler; OS threads are managed by the operating system kernel.
Memory usage: A goroutine starts with a tiny ~2 KB stack that grows automatically as needed; an OS thread typically has a fixed ~1 MB stack.
Creation cost: Goroutines are very cheap to create — you can spin up millions; OS threads are expensive, so the number you can create is limited.
Scheduling: Goroutines use the Go scheduler's M:N model (many goroutines mapped onto fewer OS threads); OS threads use the kernel's 1:1 model.
Context switching: Switching between goroutines is fast since it doesn't involve the kernel; switching between OS threads is slower because it requires a kernel context switch.
Blocking behavior: When a goroutine blocks, it parks and its underlying thread gets reused for other work; when an OS thread blocks, it blocks completely.
Communication: Goroutines typically communicate via channels; OS threads communicate using mutexes, semaphores, and other OS-level primitives.
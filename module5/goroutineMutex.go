 package module5

import (
    "fmt"
    "sync"
    "time"
)

func GoroutineMutex() {
    var mu sync.Mutex
    var counter int

    go func() {
        mu.Lock()
        counter++
        mu.Unlock()
    }()

    time.Sleep(2 * time.Second)
    fmt.Println(counter)
}

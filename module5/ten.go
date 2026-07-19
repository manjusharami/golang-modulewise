package module5

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func AtomicNumber() {
	var counter atomic.Int64
	var normalCounter int64
	var wg2 sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			for j := 1; j <= 10; j++ {
				counter.Add(1)
				normalCounter++
			}
		}()
	}

	wg2.Wait()

	fmt.Println("The value of counter is", counter.Load())
	fmt.Println("The value of normalCounter is", normalCounter)
}

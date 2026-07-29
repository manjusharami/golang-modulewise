package module5

import (
	"fmt"
	"sync"
)

func Goroutine10() {
	var wg sync.WaitGroup
	fmt.Println("execute three ")

	for i := 1; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i)
		}()
	}
	wg.Wait()
}

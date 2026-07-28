package module5

import (
	"fmt"
	"sync"
)

func GoRoutine2() {
	var wg sync.WaitGroup
	result := make([]int, 6)

	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			result[n] = n
		}(i)
	}
	wg.Wait()

	for v := range result {
		fmt.Println(v)
	}

}

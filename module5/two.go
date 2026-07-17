package module5

import (
	"fmt"
	"time"
)

func two() {
	fmt.Println("Inside method two of module5")

	// If we want to print 1 to N using go routine
	// make sure it gets printed in order
    printOneToN(5)

	time.Sleep(1 * time.Second)
}

// func printOneToN(n int) {
//   ch := make(chan int,n)
//   ch <- 1
//   for{
// 	val := <-ch
// 	fmt.Println(val)
// 	next := val + 1
// 	if next > n {
// 		close(ch)
// 		break
// 	}
// 	ch <-next
//   }
// }

func printOneToN(n int) {
	ch := make(chan int, n)

	go func() {
		for i := 0; i <= n; i++ {
			ch <- i
		}
		close(ch)
	}()

	for v := range ch {
		fmt.Println(v)
	}
}

// [1]
// Print 1 and store 1 + 1
// Print 2 and store 2 + 1
// Print 3 and store 3 + 1
// Print 4 and store 4 + 1
// Print 5 and store 5 + 1, But since 5 + 1 is greater than 5 ,user close the loop

package module5

import (
	"fmt"
	"time"
)

func hello(i int) {
	fmt.Println("print the i", i)
}
func GoRoutinetest1() {
	for i := 0; i < 5; i++ {
		go hello(i)
	}

	for i := 0; i < 3; i++ {
		go func() {
			fmt.Println("i am anomony function", i)
		}()
	}
	time.Sleep(2 * time.Second)
	// create the chaneel
	ch := make(chan int, 5)

	ch <- 10
	ch <- 11
	ch <- 13
	ch <- 14

	close(ch)

	for  v := range ch {
		fmt.Println(v)
	}

}

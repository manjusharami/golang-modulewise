package module5

import (
	"fmt"
	"time"
)

func SelectStatement() {
	fmt.Println("I am in six")

	ch1, ch2 := make(chan int), make(chan int)

	go func() {
		ch2 <- 10
	}()

	go func() {
		ch1 <- 20
	}()

	select {
	case data := <-ch1:
		fmt.Println("received data from channel 1", data)
	case data := <-ch2:
		fmt.Println("Received data from channel 2", data)
	default:
		fmt.Println("No data received yet")
	}

	time.Sleep(2 * time.Second)
}

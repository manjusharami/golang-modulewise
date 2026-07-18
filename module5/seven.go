package module5

import (
	"fmt"
	"time"
)

func Seven() {
	ch1, ch2 := make(chan int), make(chan int)

	go func() {
		for {
			ch1 <- 10
		}
	}()

	go func() {
		for {
			ch2 <- 20
		}
	}()

	for {
		select {
		case data := <-ch1:
			fmt.Println(data)
		case data := <-ch2:
			fmt.Println(data)
		}
		time.Sleep(2 * time.Second)
	}
}

package module5

import (
	"fmt"
	"time"
)

//Producer

func Gochan2() {
	ch := make(chan int)

	go Producer(ch)
	 for v := range ch {
		fmt.Println(v)
	 }
	 
}

func Producer(ch chan int) {
	for i := range 5 {
		ch <- i
	}
	   close(ch) 
	time.Sleep(3 * time.Second)
}

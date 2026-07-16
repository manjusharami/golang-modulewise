package module5

import (
	"fmt"
	"time"
)

func callMe() {
	fmt.Println("line Number 6")
}

func SimpleGoRoutine() {
	fmt.Println("simple")
	go callMe()

	go func() {
		fmt.Println("Line number 17")
	}()
	time.Sleep(1 * time.Second)
}

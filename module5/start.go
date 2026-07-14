package module5

import (
	"fmt"
	"time"
)

func hello() {
	fmt.Println("Simple hello ")
}

func Start() {
	fmt.Println("welcome to module5")

	go hello()
	time.Sleep(1 * time.Second)
}


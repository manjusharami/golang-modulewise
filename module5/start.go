package module5

import (
	"fmt"
	"time"
)

func exampleOne() {
	fmt.Println("Simple hello from examle one")
}

func Start() {
	fmt.Println("welcome to module5")

	go exampleOne()
	time.Sleep(1 * time.Second)
}

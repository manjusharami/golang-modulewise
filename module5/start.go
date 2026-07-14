package module5

import (
	"fmt"
	"time"
)

func exampleOne() {
	fmt.Println("Simple hello from examle one")
}

// func exampleTwo() {
// 	fmt.Println("Simple hello from example two")
// }

func Start() {
	go fmt.Println("welcome to module5")

	go exampleOne()
	time.Sleep(1 * time.Second)
}

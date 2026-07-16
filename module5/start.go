package module5

import (
	"fmt"
)

func exampleOne() {
	fmt.Println("Simple hello from examle one")
}

func exampleTwo() {
	fmt.Println("Simple hello from example two")

	var i = 10
	go func() {
		fmt.Println(i)
	}()
}

func Start() {
	// go fmt.Println("welcome to module5")

	// go exampleOne()

	// exampleTwo()

	// time.Sleep(1 * time.Second)

	two()
	
	//SimpleGoRoutine()

	// GoChannels()
}

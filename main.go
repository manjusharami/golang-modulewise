package main

import (
	"fmt"
	"manju/module1"
	"manju/module2"
	"manju/module3"
)

func main() {
	fmt.Println("Welcome to go")

	module1.Start()

	module2.Start()
	module3.Start()
}

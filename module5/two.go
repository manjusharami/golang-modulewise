package module5

import (
	"fmt"
	"time"
)

func two() {
	fmt.Println("Inside method two of module5")

	// If we want to print 1 to N using go routine
	// make sure it gets printed in order
	printOneToN(5)

	time.Sleep(1 * time.Second)
}

func printOneToN(n int) {
	for i := 1; i <= n; i++ {
		go fmt.Println(i)
	}
}

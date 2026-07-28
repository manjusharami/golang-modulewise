package module5

import (
	"fmt"
	"sync"
)

func Gochannel() {
	var wg sync.WaitGroup
	fmt.Println("welcome to go rountine")
	wg.Add(2)

	go runner1(&wg)
	go runner2(&wg)
	wg.Wait()

}

func runner1(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("I am runner1")
}

func runner2(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("I am runner2")
}

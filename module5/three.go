// You can edit this code!
// Click here and start typing.
package module5

import (
	"fmt"
	"sync"
)

func read(ch chan int, wg *sync.WaitGroup) {
	//defer the decrement operation of the waitgroup
	defer wg.Done()
	for data := range ch {
		fmt.Println(data)
	}
}

func WaitGroupTest() {
	var wg sync.WaitGroup

	ch := make(chan int, 10)
	// send the data to the channel
	for i := 0; i < 10; i++ {
		ch <- i
	}
	//close the buffered channel
	close(ch)
	//increment wait group counter
	wg.Add(1)
	//call the read menthod using gooutine
	go read(ch, &wg)
	//wait for counter to be decrement
	wg.Wait()
}

//Create a buffered channel of size 10 send data from main routine and read from another routine

// Main method algorithm
//Send the Data to the channel
// Close the buffered channel
// increment the waitGroup counter
//call the reader menthod using goroutine
//wait for the counter to be decremented

// Read Menthod Algorthim
// defer the decrement operation of

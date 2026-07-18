// You can edit this code!
// Click here and start typing.
package module5

import (
	"fmt"
	"sync"
)

 

func WaitGroupTest() {
	 var wg sync.WaitGroup
	 ch := make(chan int,10)
// write the channel
	 for i :=0 ;i < 10;i++{
		ch<-i

	 }

	//close
	close(ch)
	wg.Add(1)
	go read(ch,&wg)
	 
	wg.Wait()

}
 
func read(ch chan int, wg *sync.WaitGroup) {
	//defer the decrement operation of the waitgroup
	defer wg.Done()
	for data := range ch {
		fmt.Println(data)
	}
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

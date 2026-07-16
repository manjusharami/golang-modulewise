package module5

import "fmt"

func GoChannels() {
	fmt.Println("Go Channels")

	bufferedChannel()

	unBufferedChannel()
}

// When we want only limited capacity data to be written
// and read, we use bufferend channel with a capacity
func bufferedChannel() {
	myChannel := make(chan int, 3)

	myChannel <- 20
	myChannel <- 40
	myChannel <- 50

	close(myChannel)

	for data := range myChannel {
		fmt.Println(data)
	}
}

// When we don't known total number of data to be written
// and read before creating the channel, use unbufferend channel
func unBufferedChannel() {

}

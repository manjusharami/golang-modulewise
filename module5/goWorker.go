package module5

import (
	"fmt"
	"strconv"
	"time"
)

func Worker(id int) {

	fmt.Println("worker" + strconv.Itoa(id) + "started")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("worker" + strconv.Itoa(id) + "started")
}

func GoWorkerNode() {
	go Worker(1)
	go Worker(2)
	time.Sleep(2 * time.Second)

}

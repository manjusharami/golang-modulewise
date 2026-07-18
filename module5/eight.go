package module5

import (
	"fmt"
	"strconv"
	"sync"
)

var (
	wg        sync.WaitGroup
	lock      sync.Mutex
	condition = sync.NewCond(&lock)
	myNumber  = 9
	flag      = true
	ch        = make(chan string)
)

func odd() {
	defer wg.Done()
	for i := 1; i <= myNumber; i += 2 {
		lock.Lock()
		for !flag {
			condition.Wait()
		}
		ch <- "odd  data " + strconv.Itoa(i)
		flag = false
		condition.Signal()
		lock.Unlock()
	}
}

func even() {
	defer wg.Done()
	for i := 2; i <= myNumber; i += 2 {
		lock.Lock()
		for flag {
			condition.Wait()
		}
		ch <- "even data " + strconv.Itoa(i)
		flag = true
		condition.Signal()
		lock.Unlock()
	}
}

func Eight() {
	wg.Add(2)

	go odd()
	go even()

	go func() {
		wg.Wait()
		close(ch)
	}()

	for data := range ch {
		fmt.Println(data)
	}
}

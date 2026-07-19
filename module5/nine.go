package module5

import (
	"fmt"
	"sync"
	"time"
)

var (
	numWorkers = 3
	numJobs    = 10
	wg1        sync.WaitGroup
)

func worker(i int, jobs <-chan int) {
	defer wg1.Done()

	for job := range jobs {
		fmt.Printf("worker %v, jobs %v\n", i, job)
		time.Sleep(2 * time.Second)
	}
}

func WorkerPool() {
	fmt.Println("Worker Pool")

	jobs := make(chan int)

	for i := 1; i <= numWorkers; i++ {
		wg1.Add(1)
		go worker(i, jobs)
	}

	for i := 1; i <= numJobs; i++ {
		jobs <- i
	}

	close(jobs)
	wg1.Wait()

	fmt.Println("Processed all the Jobs")
}

// Make use of 3 workers ang get the avaible data of 10 in this case

package module5

import (
	"context"
	"fmt"
	"time"
)

func contextUsage() {
	fmt.Println("Context Usage")
	//contextWithTimeout()
	contextWithCancel()
}

func contextWithTimeout() {
	ctxTimeout := 5 * time.Second
	taskTimeout := 2 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	select {
	case <-time.After(taskTimeout):
		fmt.Println("In progress")
	case <-ctx.Done():
		fmt.Println("Context done")
	}
}

func contextWithCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	ctxTimeout := 5 * time.Second
	taskTimeout := 16 * time.Second
	mainWaitTimeout := 2 * time.Second
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("End the program")
				return
			default:
				fmt.Println("In Progress")
				time.Sleep(ctxTimeout)
			}
		}
	}()

	time.Sleep(taskTimeout)
	cancel()

	time.Sleep(mainWaitTimeout)
}

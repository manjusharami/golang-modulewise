package module5

import (
	"context"
	"fmt"
	"time"
)

func contextUsage() {
	fmt.Println("Context Usage")

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

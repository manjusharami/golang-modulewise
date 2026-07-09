package module4

import (
	"errors"
	"fmt"
	"math"
)

type MyError struct {
	Err string
}

// custom error handler, returning error
func (e MyError) LogError() error {
	//fmt.Println("Error Happened at" + time.Now().String() + "with the error" + e.Err)
	return fmt.Errorf(e.Err)
}

func Start() {
	fmt.Println("I am in  module 4")

	a, b := 10, 0
	var err MyError

	if b == 0 {
		err = MyError{Err: "Divide by zero is not allowed"}
	} else if b > 1000 {
		err = MyError{Err: "Larger b not allowed"}
	} else {
		fmt.Println(a / b)
	}

	fmt.Println(err.LogError())

	//sentinelError()

	panicMethod()
}

// sentinel error are  reusable predefined errrors that colors can compare
func sentinelError() {
	var errFound = errors.New("Record not found")
	a := []int{2, 3, 4, 5}
	fmt.Println("sentinelError")
	if len(a) > 6 {
		fmt.Println(errFound)
	}
}

func panicMethod() {
	a := 1000.0
	if a < 0 {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovering from panic")
			}
		}()
		panic("Negative nummber is not allowed")
	}
	fmt.Println(math.Sqrt(a))
}

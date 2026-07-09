package module4

import (
	"fmt"
	"time"
)

type MyError struct {
	Err string
}
 // custom error handlers 
func (e MyError) LogError() {
	fmt.Println("Error Happened at" + time.Now().String() + "with the error" + e.Err)
}

func Start() {
	fmt.Println("I am in  module 4")

	a, b := 10, 2000
	var err MyError

	if b == 0 {
		err = MyError{Err: "Divide by zero is not allowed"}
	} else if b > 1000 {
		err = MyError{Err: "Larger b not allowed"}
	} else {
		fmt.Println(a / b)
	}

	err.LogError()
}

package module4

import (
	"errors"
	"fmt"
	"math"
	"os"
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

	//panicMethod()

	readingFile()

	writeFile()

	appendFile()

	removeFromFile()
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

func readingFile() {
	data, err := os.ReadFile("data/manju.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func writeFile() {
	inputData := []byte("i am fifa world cup messi")
	err := os.WriteFile("data/messi.txt", inputData, 0644)
	if err != nil {
		panic("could not written into the file")
	}
	fmt.Println("File Created Succesful")
}

func appendFile() {
	file, err := os.OpenFile("data/appendExample.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic("Cant able to open the file")
	}
	defer file.Close()
	file.WriteString("\nthis new content appended")
}

func removeFromFile() {
	// remove the last line from the text file and call the method
}

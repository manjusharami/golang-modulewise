package pratice

import "fmt"

func Solution() {
	A := []int{1, 2, -1, 3, 3}

	count := 0
	index := 0

	for index != -1 {
		count++
		index = A[index]
		fmt.Println("index",index)
		// fmt.Println("A[index]",A[index])

	}
	fmt.Println(count)

}

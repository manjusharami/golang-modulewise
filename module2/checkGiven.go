//Check if a given target value exists in the array.

package module2

import "fmt"

func checkGiven() {

	array1 := []int{1, 3, 5, 3, 10, 30}
	target := 100
	hasPresent := false

	for _, v := range array1 {
		if v == target {
			hasPresent = true
		}
	}

	if hasPresent{
		fmt.Println("The Target is Present")
	}else{
		fmt.Println("the Target is not Present")
	}

}

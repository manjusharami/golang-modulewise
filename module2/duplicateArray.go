package module2

import "fmt"

func duplicateArray() {

	array1 := []int{1,2, 4, 5, 7}
	seen := make(map[int]bool)
	hasDuplicate := false
	for _, v := range array1 {
		if seen[v] {
			 hasDuplicate =true
			 break	
		}
		seen[v] = true
	}
	if hasDuplicate {
		fmt.Println("The array contail duplicate values")
	} else {
		fmt.Println("The array does not contail duplicate values")
	}

}

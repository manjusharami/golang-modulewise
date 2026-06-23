// Determine whether two arrays have at least one common element.
package module2

import "fmt"

func commonElement() {
	array1 := []int{9, 7, 5, 5, 6}
	array2 := []int{9, 5, 7, 4, 3, 4}
	commonEle := []int{}
	seen := make(map[int]bool)
	for _, v1 := range array1 {
		seen[v1] = true
	}
	for _, v2 := range array2 {
		if seen[v2] {
			commonEle = append(commonEle, v2)

		}

	}
	if len(commonEle) >= 0 {
		fmt.Println("The Array  have common element", commonEle)
	} else {
		fmt.Println("Array doest not have common elements", commonEle)
	}

}

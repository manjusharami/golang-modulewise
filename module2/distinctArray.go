//1. Find the number of distinct elements in an array.//

package module2

import "fmt"

func distinctArray() {
	array1 := []int{1, 3, 3, 3, 5, 2, 6, 8, 9, 8}
	seen := make(map[int]bool)
	unique := []int{}

	for _, v := range array1 {
		 
		if !seen[v] {
			unique = append(unique, v)
			seen[v] = true
		}

	}
	fmt.Println(unique)
	fmt.Println(len(unique))
}

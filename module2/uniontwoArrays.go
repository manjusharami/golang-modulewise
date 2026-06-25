//Union of Two Arrays

package module2

import (
	"fmt"
	"sort"
)

func unionofTwoArrays() {
	array1 := []int{2, 4, 2, 2, 4, 2, 100, 300}
	array2 := []int{2, 3, 4, 5, 6, 7, 8}

	set := make(map[int]struct{})
	result := []int{}

	for _, v := range array1 {

		if _, exists := set[v]; !exists {
			set[v] = struct{}{}
			result = append(result, v)
		}

	}

	for _, v := range array2 {

		if _, exists := set[v]; !exists {
			set[v] = struct{}{}
			result = append(result, v)
		}

	}
	sort.Ints(result)

	fmt.Println(result )
}

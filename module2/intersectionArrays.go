//Intersection of Two Arrays
package module2

import (
	"fmt"
)

func intersectionArray() {
	array1 := []int{1, 4, 3, 2, 2, 4, 2}
	array2 := []int{10,3, 4, 4, 2, 2}
	set := make(map[int]struct{})
	result := []int{}

	for _, arr1 := range array1 {
		set[arr1] = struct{}{}
	}
	for _, arry2 := range array2 {
		if _, exits := set[arry2]; exits {
			result = append(result, arry2)
			delete(set,arry2)
		}
	}
	fmt.Println(result)
}

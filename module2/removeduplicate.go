// Remove duplicate elements from an array and print unique values using hashMAP.
package module2

import "fmt"

func removeDuplicate() {
	array1 := []int{1, 2, 3, 4, 4, 5, 2, 2,8,9}
	removeDuplicateEle := []int{}
	seen := make(map[int]int)

	for _, v := range array1 {
		seen[v]++
		if seen[v] == 1 {
			removeDuplicateEle = append(removeDuplicateEle, v)
		}
	}
	fmt.Println(removeDuplicateEle)

}

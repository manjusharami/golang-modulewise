//2. Remove Duplicates from an Array using HashSet

package module2

import "fmt"

func removeDuplicatesSet() {
	array := []int{10, 20, 20, 40, 50}
	setSeen := make(map[int]struct{})
	result := []int{}

	for _, v := range array {
		if _, exits := setSeen[v]; !exits {
			setSeen[v] = struct{}{}
			result = append(result, v)

		}

	}
	fmt.Println(result)

}

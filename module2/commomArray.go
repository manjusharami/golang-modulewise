//Check if Two Arrays Have Common Elements

//step -create the hashset
//using range and for loop put all element. of array 1 in hashmap
// using ranage and for loop check  element in arrays is in hashmap if it is there then  add it to result

package module2

import "fmt"

func commonArray() {
	array1 := []int{10, 20, 30, 40, 20, 10}
	array2 := []int{10, 20, 3, 4, 5, 8, 9, 30, 10}
	hashSet := make(map[int]struct{})
	resultSet := make(map[int]struct{})
	result := []int{}

	for _, v := range array1 {
		hashSet[v] = struct{}{}
	}
	for _, v := range array2 {
		if _, ok := hashSet[v]; ok {
			if _, seen := resultSet[v]; !seen {
				resultSet[v] = struct{}{}
				result = append(result, v)
			}
		}
	}
	fmt.Println("the result of common element in the array is", result)
}

// // create the hashset
// add element to hashmap if not added already
// for every 3 element find the distinct and count it
// Count Distinct Elements in Every Window of Size K
package module2

import "fmt"

func distinctElements() {

	array1 := []int{1, 40, 20, 10, 20, 20, 1, 9, 4, 10}
	hashMAp := make(map[int]int)
	k := 4
	//count := 0
	//result := []int{}

	for _, v := range array1[:k] {

		hashMAp[v]++
	}
	fmt.Println(hashMAp)
}

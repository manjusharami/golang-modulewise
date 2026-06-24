// Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.

// You may assume that each input would have exactly one solution, and you may not use the same element twice.

// You can return the answer in any order.
//Count the number of pairs with a given sum using a HashSet.

package module2

import "fmt"

func countPairsWithSum() {
	array1 := []int{1, 2, 4, 5, 5, 6, 7, 5, 5}
	target := 10
	seen := make(map[int]bool)
	result := []int{}
	count := 0
	for i, v := range array1 {
		need := target - v
		if seen[need] {
			result = append(result, i) //retutn the index of the slice
			count++

		}
		seen[v] = true
	}
	fmt.Println(result)
	fmt.Println(count)
}

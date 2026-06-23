//Find all elements that appear exactly once in an array.

package module2

import "fmt"

func exactlyOnce(){
	array1 := []int {1,3,4,5,4,3,3,2}
	seen := make(map[int]int)
	exactlyOnceele := []int{}

	for _,v := range array1{
		seen[v]++
		if seen[v]==1{
			exactlyOnceele = append(exactlyOnceele,v)
		}
	}
	fmt.Println("The element that apper exactly once in an array",exactlyOnceele)
}
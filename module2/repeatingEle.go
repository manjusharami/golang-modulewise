package module2

import "fmt"

func repeatingElement() {
	array1 := []int{8, 9, 7, 6, 4, 6, 8, 6, 1}
	repeatingEle := []int{}
	seen := make(map[int]int)
	for _, v := range array1 {
		seen[v]++
		if seen[v] == 2 {
			repeatingEle = append(repeatingEle, v)
		}
	}
	fmt.Println(repeatingEle)
}

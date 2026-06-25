//create the hashmap   int bool  put the element in the map and make ture and addend to result

package module2

import "fmt"

func distinctArray(){
	array1 :=[]int{10,20,30,9,20,20,20}
	hashMap :=make(map[int]bool)
	result := []int{}

	for _ , v := range array1{
		if _,ok := hashMap[v];!ok{
			hashMap[v] = true
			result = append(result, v)
		}
	}
	fmt.Println(result)
}
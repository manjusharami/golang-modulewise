package module2

import "fmt"

func uniqueElemt(){
	array1 := []int{1,8,9,8,7,7,6,8}
	seen :=make(map[int]int)
	uniqueElemt := []int{}
	for _,v := range array1 {
		seen[v]++
	 
	}
	for v,count := range seen{
		if count ==1{
			uniqueElemt = append(uniqueElemt,v)
		}
	}

	fmt.Println("the unique element in the array is",uniqueElemt)
	fmt.Println(" unique element in the array is",len(uniqueElemt))

}
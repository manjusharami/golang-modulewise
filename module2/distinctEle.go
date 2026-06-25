//Find All Distinct Elements in an Array

//Step1 create variable  set := map([]int struct{}) and slice
//Step 2   create for loop and put all the values in the array to set and  appened to result the value

package module2

import "fmt"

func distincEle(){
	array1 :=[]int{10,20,60,20,1,2,4,5,10,20,60}
	set :=make(map[int]struct{})
	result := []int{}
	for _,v := range array1 {
		if _,ok := set[v];!ok{
			set[v]=struct{}{}
			result = append(result,v )

		}
	}
	fmt.Println("the array is",result)
}

//1. Find the number of distinct elements in an array.//

// package module2

// import "fmt"

// func distinctArray() {
// 	array1 := []int{1, 3, 3, 3, 5, 2, 6, 8, 9, 8}
// 	seen := make(map[int]bool)
// 	unique := []int{}

// 	for _, v := range array1 {
		 
// 		if !seen[v] {
// 			unique = append(unique, v)
// 			seen[v] = true
// 		}

// 	}
// 	fmt.Println(unique)
// 	fmt.Println(len(unique))
// }

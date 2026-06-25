//Find Pair with Given Sum using hashSet
//Get array and target sum
//step 1 creat hashset ,
///setp 2 calulate need= target -get element in the array
//check needed in the array if the needed is present appened to result the index

package module2

import "fmt"

func pairSum(){
	array1 := []int{100,200,900,100,1,4,90}
	targetSum := 1000
	hashSet := make(map[int]struct{})
	result :=[]int {}

	for i,v := range array1 {
		need := targetSum-v
		hashSet[v]=struct{}{}
		if _,ok := hashSet[need] ;ok{
			result =append(result,i)
		}
	 

	}
	fmt.Println(result)
}

//1000-100=900 check the number present in the array ,-if yes  add to result 1=itreation
//1000-200=800-not there in array-2 itreation
//1000-900-100 is present in array add it
//1000-100=900-is presnt in array alreaded add ignore
//1000-1 - 9999 is presnt in arrya no 
//1000-4=996 no
//1000-90=890


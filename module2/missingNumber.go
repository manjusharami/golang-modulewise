//Find Missing Number

package module2

import "fmt"

func missingNumber(){
	nums := []int{1,3,4,5,3}
	set := make(map[int]struct{})
	
	
	for _,v := range nums{
		set[v]=struct{}{}
	}
	for i:=1 ;i<len(nums);i++{
		 if _ ,exits := set[i];!exits{
			fmt.Println("the missing number is ",i)
		 }
	}
}
//

package pratice

import "fmt"

func FindArray() {
	array1 := []int{1, 2, 3, 4}
	n := len(array1)
	result := make([]int, n)
	 
	prefix :=1
	 for i:=0 ; i<n ;i++ {
		result[i] = prefix
		fmt.Println(prefix)
		prefix *=array1[i]
		fmt.Println("prefix",result)
	 }
	 sufix := 1
	 for i :=n-1 ;i>=0 ;i--{
		 
		result[i]*=sufix
		fmt.Println("sufix",sufix)
		fmt.Println(array1[i])
		sufix*=array1[i]
		fmt.Println("update suffix",sufix)
		fmt.Println("suffix",result)
	 }
	fmt.Println("result",result)
}

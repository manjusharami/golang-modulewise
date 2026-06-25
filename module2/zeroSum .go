//Number of Zero-Filled Subarrays
// create hashset
//loop with for loop check if 9 present in the array if yes add the count++

package module2

import "fmt"

func zeroSum() {
	array1 := []int{10, 0, 1, 0, 3, 4, 5, 0, 0, 1, 0, 0}

	count := 0
	subArray := 0

	for _, v := range array1 {
		if v == 0 {
			subArray++
			count += subArray
		} else {
			subArray = 0
		}

	}
	fmt.Println("the subarray is zero is",count)

}
// First Iteration subArray =0 
//2 Iteration subArray=1 count=1
//3 Iteration subArray =0 
//4 Iteration Subarray=1 count=2
//5,6,7  iteration subarray=0 count=2
//8 iteration subaray=1 count=3
//9 iteration subarrayy=2 count=5
//10 iteration subarrray=0
//11 iteration subarray=1 count=6
//12 iteration subarray =2 count =8
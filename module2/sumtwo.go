package module2

import "fmt"

 


func sumTwo(nums[]int ,target int) []int{

	present := make(map[int]int)
	result := []int{}

	for i,v := range nums{
		need :=target-v

	if idx,ok := present[need];ok {
		 result = append(result,idx,i)
		 return  result
	}
  present[v]=i
}
return result
}

func findsumtarget()  {
	nums := []int{1,2,3,4,6}
	tar := 9
	result := sumTwo(nums,tar)
    fmt.Println(result)

}

// Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.

// You may assume that each input would have exactly one solution, and you may not use the same element twice.

// You can return the answer in any order

// Example 1:

// Input: nums = [2,7,11,15], target = 9
// Output: [0,1]
// Explanation: Because nums[0] + nums[1] == 9, we return [0, 1].
// Example 2:

// Input: nums = [3,2,4], target = 6
// Output: [1,2]
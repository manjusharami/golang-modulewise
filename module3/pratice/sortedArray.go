// . Search in Rotated Sorted Array-

 package pratice

import "fmt"

func search(nums []int, target int) int {
    left, right := 0, len(nums)-1

    for left <= right {
        mid := (left + right) / 2
		fmt.Println("midmidmid",mid)
        if nums[mid] == target {
            return mid
        }
		fmt.Println("nums[left]",nums[left])
		fmt.Println("nums[left]",nums[mid])
        // Left half is sorted
        if nums[left] <= nums[mid] {
            if target >= nums[left] && target < nums[mid] {
                right = mid - 1
            } else {
                left = mid + 1
            }
        } else {
            // Right half is sorted
            if target > nums[mid] && target <= nums[right] {
                left = mid + 1
            } else {
                right = mid - 1
            }
        }
    }

    return -1
}

func SortedArray() {
    nums := []int{4, 5, 6, 7, 0, 1, 2}
    target := 0

    result := search(nums, target)   //   call binary search

    fmt.Println(result)
}

// . Median of Two Sorted Arrays

package pratice

import (
	"fmt"
)

func SolutionMedium() {
	a1 := []int{10, 20, 30, 40, 50}
	a2 := []int{60, 70, 80}
	result := []int{}
	i := 0
	j := 0

	for i < len(a1) && j < len(a2) {
		if a1[i] < a2[j] {
			result = append(result, a1[i])
			i++
		} else  {
			result = append(result, a2[j])
			j++
		}
		fmt.Println("result", result)

		for i < len(a1) {
			result = append(result, a1[i])
			i++
		}

		for j < len(a2) {
			result = append(result, a2[j])
			j++
		}
		n := len(result)

		if n%2 == 1 {
			mediam :=result[n/2]
			fmt.Println("The Medium of the array is", mediam)
		} else {
			mediam :=float64(result[n/2-1]+result[n/2])/2
			fmt.Println("the Medium of the Array",mediam)
		}

	}
}

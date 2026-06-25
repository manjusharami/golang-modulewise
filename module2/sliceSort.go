// You can edit this code!
// Click here and start typing.
package module2

import (
	"fmt"
	"sort"
)

func sliceSort() {
	slice1 := []int{3, 5, 6, 7, 4, 64}
	sort.Slice(slice1, func(i, j int) bool {
		return slice1[i] < slice1[j]
	})
	fmt.Println(slice1)
	slice2 := []string{"manju", "chinnu"}
	sort.Slice(slice2, func(i, j int) bool {
		return slice2[i] < slice2[j]
	})
	fmt.Println(slice2)
}
package module2

import (
	"fmt"
	"maps"
	"slices"
)

func normalSliceCopy() {
	mySlice1 := []int{1, 3, 4, 4, 2}
	var mySlice2 []int
	for _, v := range mySlice1 {
		mySlice2 = append(mySlice2, v)
	}

	fmt.Println(mySlice1, mySlice2)
	mySlice2[2]=56666
	fmt.Println(mySlice1,mySlice2)
}

func efficientSliceCopy() {
	mySlice1 := []string{"manju", "sweety", "chinnu"}

	mySlice2 := slices.Clone(mySlice1)
	fmt.Println(mySlice2)

	mySlice2[1]="dhajshdhjashgdjha"
	fmt.Println(mySlice1,mySlice2)
}

func mapCopy() {
	map1 := make(map[int]int)
	map1[0] = 100
	map1[2] = 300
	map1[3] = 200

	map2 := maps.Clone(map1)
	fmt.Println(map2)

	map2[0] = 1234
	fmt.Println(map1[0], map2[0])
}

func copyFun() {
	normalSliceCopy()
	efficientSliceCopy()
	mapCopy()
}

package module2

import "fmt"

func setOperations() {
	fmt.Println("inside hash set")

	// Hashset creation
	myHashSet := make(map[int]struct{})

	// Hashset insert
	myHashSet[100] = struct{}{}
	myHashSet[300] = struct{}{}
	myHashSet[500] = struct{}{}
	myHashSet[200] = struct{}{}
	myHashSet[400] = struct{}{}

	// fmt.Println(myHashSet)

	// Checking if the value exist in the HashSet
	// value, present := myHashSet[100]

	// fmt.Println(myHashSet, value, present)

	// value1, present1 := myHashSet[1001]

	// fmt.Println(myHashSet, value1, present1)

	fmt.Println(myHashSet)
	myHashSet[200] = struct{}{} // this data will not stored because it is already present
	fmt.Println(myHashSet)

	delete(myHashSet, 400)

	fmt.Println(myHashSet)

	// itreate over the hashset

	for key := range myHashSet {
		fmt.Println(key)
	}

}

//use hashSet for quick search
// and/or finding unquie elements from large dataset

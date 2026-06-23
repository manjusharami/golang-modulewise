package module2

import (
	"fmt"
)

func mapOperation() {
	fmt.Println("inside hash map")

	//HashMap Creation
	myHashMap := make(map[int]string)

	myHashMap[1] = "manjusha"
	myHashMap[4] = "sweety"
	myHashMap[2] = "chinnu"
	myHashMap[3] = "papa"
	fmt.Println(myHashMap) // full map will printed

	// define a search key
	userInput1 := 2
	value1, present1 := myHashMap[userInput1]

	userInput2 := 300
	value2, present2 := myHashMap[userInput2]

	fmt.Println("24", value1, present1)

	fmt.Println("26", value2, present2)

	myHashMap[4] = "aruna"
	fmt.Println(myHashMap) // full map will be printed with replace with sweety with aruna
	delete(myHashMap, 4)   // delete sweety

	for key, value := range myHashMap {
		fmt.Println(key, value) //manjusha,
	}

}

//1 manjusha
//2 chinnu
//3 papa

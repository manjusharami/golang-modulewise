package module2

import "fmt"

func mapOperation() {
	fmt.Println("inside hash map")

	//HashMap Creation
	myHashMap := make(map[int]string)

	myHashMap[1] = "manjusha"
	myHashMap[4] = "sweety"
	myHashMap[2] = "chinnu"
	myHashMap[3] = "papa"

	fmt.Println(myHashMap)
}

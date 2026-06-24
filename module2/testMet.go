
// Click here and start typing.
package module2

import (
	"fmt"
)

func testMet()  {
	 var result int
	defer func() {
		result = 20
		fmt.Println("12", result)
	}()

	defer func() {
		result = 2020
		fmt.Println("17", result)
	}()
	fmt.Println("18", result)
}

 

//17  2020
// 12 20
//18 20

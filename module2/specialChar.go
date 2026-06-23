// You can edit this code!
// Click here and start typing.
package module2

import (
	"fmt"
	"reflect"
)

func SpecialChars() {
	myEmoji1 := '😀'
	fmt.Println(myEmoji1, reflect.TypeOf(myEmoji1))

	myEmoji2 := "😀"
	fmt.Println(myEmoji2, reflect.TypeOf(myEmoji2))

	myString := "Hello"

	for index, value := range myString {
		fmt.Println(index, string(value), reflect.TypeOf(value))
	}

}

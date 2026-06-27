// Interview Question
// https://medium.com/@praffulmishra/struct-alignment-and-padding-in-go-a-comprehensive-guide-ae928d5a9d5e
package pratice

import (
	"fmt"
	. "unsafe"
)

type A struct {
	p int64
	q int16
	r int8
}
type B struct {
	p int16
	q int64
	r int8
}

func structMem() {
	obj1 := A{}
	obj2 := B{}
	fmt.Println(Sizeof(obj1), Sizeof(obj2))
}

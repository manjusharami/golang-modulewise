package module4

import "fmt"

type Person struct {
	Name string 
	Id   int
}
type Animal struct {
	Name string
	Id   int
}
type Walking interface {
	Walk() string
}

func (p Person) Walk() string {
	return  p.Name
}

func (p Animal) Walk() string {
	return  p.Name
}
func TestStrcut() {
w := []Walking{
		Person{Name: "MAJU", Id: 10},
		Person{Name: "MAJU2", Id: 11},
		Animal{Name:"Happy", Id:1},
	}
	 

 
	 


	 for _, v := range w {
		fmt.Println("The Name is ", v.Walk())
	}

}

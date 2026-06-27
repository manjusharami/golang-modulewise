package module3

import "fmt"

type Empolyee struct {
	Name   string
	Salary float64
}

func StructureBasic() {
	fmt.Println("Struct Basic")

	emp1 := Empolyee{
		Name:   "manjusha",
		Salary: 100000.90,
	}
	emp2 := Empolyee{
		Name:   "Cherishma",
		Salary: 200000.09,
	}
	fmt.Printf("Salary of %v is %v\n", emp1.Name, emp1.Salary)
	fmt.Printf("Salary of %v is %v\n", emp2.Name, emp2.Salary)
}

// Salary of manjusha is 10000.90

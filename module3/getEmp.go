package module3

import (
	"fmt"
	"manju/module3/model"
)

func getEmp() {
	emp := []model.Employee{
		{Ename: "manju", EId: 2},
		{Ename: "mmmanju", EId: 3},
	}
	for _, v := range emp {
		fmt.Println("the employee name and Eid is", v.Ename, v.EId)

	}
}

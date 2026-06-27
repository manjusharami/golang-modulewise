package pratice

import (
	"fmt"
	"sort"
)

func sortEmpByname(emp []Employee) []Employee {
	sort.Slice(emp, func(i, j int) bool {
		return emp[i].Ename < emp[j].Ename
	})
	return emp
}

func sortEmpByNumber(emp []Employee) []Employee {

	sort.Slice(emp, func(i, j int) bool {
		return emp[i].EId < emp[j].EId
	})
	return emp
}

func uniqueName(emp []Employee) []Employee {
	hashMap := make(map[string]bool)
	unique := []Employee{}
	for _, v := range emp {
		if _, ok := hashMap[v.Ename]; !ok {
			hashMap[v.Ename] = true
			unique = append(unique, v)
		}
	}
	return unique
}

func GetEmp() {
	emp := []Employee{
		{Ename: "manju", EId: 2},
		{Ename: "mmmanju", EId: 3},
		{Ename: "chinnu", EId: 4},
		{Ename: "sweety", EId: 2},
		{Ename: "amu", EId: 1},
		{Ename: "amu", EId: 1},
		{Ename: "amu", EId: 1},
	}

	for _, v := range emp {
		fmt.Println("the employee name and Eid is", v.Ename, v.EId)
	}
	sortTheEmp := sortEmpByname(emp)
	fmt.Println("Sorted by name", sortTheEmp)
	fmt.Println("_____________________")
	sortTheEmpByNumber := sortEmpByNumber(emp)
	fmt.Println("Sorted by Number", sortTheEmpByNumber)
	fmt.Println("_____________________")
	fmt.Println("remove duplicate number")

	removeDup := uniqueName(emp)
	fmt.Println("remove the duplicate", removeDup)

}

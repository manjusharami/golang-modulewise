// Create two students with there own name and roll number and
// Write a method which will find out student with higher roll number and
// Another menthod which will find out student with larger name

package module3

import "fmt"

type Students struct {
	name       string
	rollNumber int
}

//Create the Constructor -Ir is used to inilized the object and create stuct and vailde the stucture

func NewStudents(ename string, rnumber int) Students {
	return Students{name: ename, rollNumber: rnumber}

}

func hingerRollNumber(student []Students) Students {
	highest := student[0]

	for _, v := range student {
		if v.rollNumber > highest.rollNumber {
			highest = v
		}
	}
	return highest
}

func longestNameOfStudent(student []Students) Students {
	highest := student[0]

	for _, v := range student {
		if len(v.name) > len(highest.name) {
			highest = v
		}
	}
	return highest
}

func (s *Students) GetName() string {
	return s.name
}

func (s *Students) GetRollNumber() int {
	return s.rollNumber
}

func (s *Students) setName(str string) {
	s.name = str
}

func (s *Students) setRollnumber(r int){
	s.rollNumber=r
}
func createStudent() {
	student1 := NewStudents("manjus", 300)
	student2 := NewStudents("sweety", 33)
	student3 := NewStudents("chinnureddy", 100)
	student4 := NewStudents("bunnny", 100)
	stu := []Students{student1, student2, student3, student4}
	student := hingerRollNumber(stu)
	fmt.Println("the highest roll number of the student is", student)
	logestName := longestNameOfStudent(stu)
	fmt.Println("logestName", logestName)

}

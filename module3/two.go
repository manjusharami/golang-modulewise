package module3

import "fmt"

type Student struct {
	Name string
	Roll int
}

// Constructor in GoLang Syntax -  This is done by convention over configuration approach
func NewStudent(str string, r int) Student { return Student{Name: str, Roll: r} }

func (s *Student) GetName() string { return s.Name }

func (s *Student) GetRoll() int { return s.Roll }

func (s *Student) SetName(str string) { s.Name = str }

func (s *Student) SetRoll(r int) { s.Roll = r }

// toString() method simulation like java or other object oriented languages
func (s *Student) ToString() string {
	return fmt.Sprintf("Student Name is %v and Roll is %v\n", s.GetName(), s.GetRoll())
}

func (s Student) GetLengthOfStudentName() int {
	return len(s.GetName())
}

func StructurePratice() {
	fmt.Println()
	// Creating Object
	s1 := NewStudent("manjusha", 2)

	// Display the Object
	fmt.Println(s1.ToString())

	// Change s1 students rollNo to 10
	s1.SetRoll(10)

	// Display the Object
	fmt.Println(s1.ToString())

	// Print only the name of s1 student
	fmt.Println(s1.GetName())

	// Find lenght of Student Name
	fmt.Println(s1.GetLengthOfStudentName())
}
// Create two students with there own name and roll number and 
// Write a method which will find out student with higher roll number and
// Another menthod which will find out student with larger name 


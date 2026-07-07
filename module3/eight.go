package module3

import "fmt"

// ======================================================

// 1. Structures in Go

// ======================================================

type Person struct {
	Name string

	Age int
}

// ======================================================

// 2. Struct Methods

// ======================================================

// Method attached to Person

func (p Person) PrintDetails() {

	fmt.Println("Name:", p.Name)
	fmt.Println("Age :", p.Age)

}

// ======================================================

// 3. Nested Structs

// ======================================================

type Address struct {
	City    string
	Country string
}

type Employee struct {
	Name string

	Address Address
}

// ======================================================

// 4. Pointers and Memory Handling

// ======================================================

// Receives pointer so original object gets modified

func updateAge(p *Person) {

	p.Age = 99

}

// ======================================================

// 5. Methods and Receivers

// ======================================================

// Value Receiver (works on copy)

func (p Person) ChangeNameValueReceiver(name string) {

	p.Name = name

}

// Pointer Receiver (works on original object)

func (p *Person) ChangeNamePointerReceiver(name string) {

	p.Name = name

}

// ======================================================

// 6. Interfaces and Polymorphism

// ======================================================

type Animal interface {
	Speak()
}

type Dog struct {
	Name string
}

func (d Dog) Speak() {

	fmt.Println(d.Name, "says Woof")

}

type Cat struct {
	Name string
}

func (c Cat) Speak() {

	fmt.Println(c.Name, "says Meow")

}

// ======================================================

// 7. Composition

// ======================================================

type Engine struct {
	HorsePower int
}

func (e Engine) Start() {

	fmt.Println("Engine Started")

}

// Car HAS-A Engine

type CarComposition struct {
	Brand string

	Engine Engine
}

// ======================================================

// 8. Embedding

// ======================================================

// Car embeds Engine

// Engine fields and methods are promoted

type CarEmbedding struct {
	Brand string

	Engine
}

// ======================================================

// Demo Methods

// ======================================================

// 1. Structures

func structuresExample() {

	fmt.Println("========== Structures ==========")

	p := Person{

		Name: "Rahul",

		Age: 25,
	}

	fmt.Println(p)

	fmt.Println()

}

// 2. Struct Methods

func structMethodsExample() {

	fmt.Println("========== Struct Methods ==========")

	p := Person{
		Name: "Amit",
		Age:  30,
	}

	p.PrintDetails()

	fmt.Println()

}

// 3. Nested and Anonymous Structs

func nestedStructExample() {

	fmt.Println("========== Nested Struct ==========")

	e := Employee{
		Name: "Rohan",
		Address: Address{
			City:    "Bangalore",
			Country: "India",
		},
	}

	fmt.Println(e.Name)
	fmt.Println(e.Address.City)
	fmt.Println(e.Address.Country)

	fmt.Println()

	fmt.Println("========== Anonymous Struct ==========")

	student := struct {
		Name string

		Marks int
	}{

		Name: "Ankit",

		Marks: 95,
	}

	fmt.Println(student)

	fmt.Println()

}

// 4. Pointers

func pointerExample() {

	fmt.Println("========== Pointers ==========")

	p := Person{

		Name: "Neha",

		Age: 20,
	}

	fmt.Println("Before:", p)

	updateAge(&p)

	fmt.Println("After :", p)

	fmt.Println()

}

// 5. Methods & Receivers

func receiverExample() {

	fmt.Println("========== Value Receiver ==========")

	p := Person{
		Name: "Karan",
		Age:  27,
	}

	p.ChangeNameValueReceiver("Changed")

	fmt.Println(p.Name)

	fmt.Println()

	fmt.Println("========== Pointer Receiver ==========")

	p.ChangeNamePointerReceiver("Updated")

	fmt.Println(p.Name)

	fmt.Println()

}

// 6. Interfaces

func interfaceExample() {

	fmt.Println("========== Interfaces ==========")

	var a Animal

	a = Dog{Name: "Bruno"}

	a.Speak()

	a = Cat{Name: "Kitty"}
	a.Speak()

	fmt.Println()

}

// 7. Composition

func compositionExample() {

	fmt.Println("========== Composition ==========")

	car := CarComposition{

		Brand: "Toyota",

		Engine: Engine{

			HorsePower: 180,
		},
	}

	fmt.Println(car.Brand)

	// Need Engine field
	fmt.Println(car.Engine.HorsePower)
	//fmt.Println(car.HorsePower)

	car.Engine.Start()

	fmt.Println()

}

// 8. Embedding

func embeddingExample() {

	fmt.Println("========== Embedding ==========")

	car := CarEmbedding{

		Brand: "Honda",

		Engine: Engine{
			HorsePower: 220,
		},
	}

	fmt.Println(car.Brand)

	// Direct access because Engine is embedded

	fmt.Println(car.HorsePower)

	car.Start()

	// These also work

	fmt.Println(car.Engine.HorsePower)

	car.Engine.Start()

	fmt.Println()

}

// ======================================================

// 9. Type Assertions

// ======================================================

func typeAssertionExample() {

	fmt.Println("========== Type Assertion ==========")

	// var x interface{} = "Hello Go"
	var x any = "Hello Go"

	str := x.(string)

	fmt.Println(str)

	fmt.Println()

}

// ======================================================

// 10. Type Switch

// ======================================================

func typeSwitchExample() {

	fmt.Println("========== Type Switch ==========")

	checkType(100)
	checkType("Go")
	checkType(true)

	checkType(10.5)

	fmt.Println()

}

func checkType(v any) {

	switch value := v.(type) {

	case int:
		fmt.Println("Integer:", value)

	case string:
		fmt.Println("String :", value)

	case bool:

		fmt.Println("Boolean:", value)

	default:

		fmt.Println("Unknown Type")

	}

}

// ======================================================

func CompleteModule3() {

	structuresExample()

	structMethodsExample()

	nestedStructExample()

	pointerExample()

	receiverExample()

	interfaceExample()

	compositionExample()

	embeddingExample()

	typeAssertionExample()

	typeSwitchExample()

}

package module3

import "fmt"

type AnimalInterface interface {
	Eat(str string)
}

type Lion struct{}

func (a Lion) Eat(str string) {

	fmt.Println("Lion Eating Every Minute")

}

type Tiger struct{}

func (a Tiger) Eat(str string) {

	fmt.Println("Tiger Eating Fish Only Everyday and", str, "is Not Allowed")

}

func StartEating(animal AnimalInterface, food string) {

	animal.Eat(food)

}

func AnimalInter() {

	l := Lion{}

	t := Tiger{}
	StartEating(l, "chicken")

	StartEating(t, "burger")

}

//Interface Makes It Easier For Users To Depend On Abstraction, Not Actual Implementation

// Lion Eating Every Minute

// Tiger Eating Fish Only Everyday and burger is Not Allowed

//interface is defines menthod declaration and it is not menthode implementation
//what is the use of the interface is loose couping
//loose couping-code need easier to change -This makes systems easier to change without breaking everything.
//Clear contracts — Interfaces define a strict set of methods a class must implement, improving consistency.
//Dependency Injection — Modern frameworks rely heavily on interfaces to inject dependencies.
//Polymorphism — You can treat different objects the same way if they implement the same interface.

package module3

import "fmt"

type Bank interface {
	Deposit(amount float64)
	Withdraw(amount float64)
	Balance() float64
}

type Hsbc struct {
	blance float64
}
 
func NewHsbc(inintailBalance float64) Bank{
	return &Hsbc{blance: inintailBalance}
}

func (b *Hsbc) Deposit(amount float64) {
	b.blance += amount
}

func (b *Hsbc) Withdraw(amount float64) {
	b.blance -= amount
}
func (b *Hsbc) Balance() float64 {
	return b.blance
}

func CallInterface() {
	 // var b Bank = &Hsbc{}

	 var b Bank =NewHsbc(900)
 
	b.Deposit(100)
	b.Withdraw(50)
	b.Balance()
	fmt.Println(b.Balance())
}

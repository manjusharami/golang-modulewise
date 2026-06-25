//Happy number

//the idea 21=sqt(2)+sqt(1)=4+1=5=sqt(5)=25

package module2

import "fmt"

func squrenumber( num int) int {
 
	// set := make(map[int]struct{})
	 
    squroot := 0
 
	for num > 0 {
		digit := num % 10
		squroot += digit * digit
		num/=10

	}

	 return squroot

}

func happyNumber(){
		num := 19
		set := make(map[int]struct{})
        
		for num != 1{
			if _,ok := set[num] ;ok {
            fmt.Println("No Happy number return")
			return
			}
			set[num] =struct{}{}
			num = squrenumber(num)
		}
	fmt.Println("Happy Number")
}
package module2

import "fmt"

 func longestSeqSet(){
    num := []int{1, 3, 5, 3, 10, 11, 12, 13, 14, 60, 1, 2, 3, 15}
	set := make(map[int]struct{})
	 
	// count :=0
	for _, v := range num{
		 
			set[v] =struct{}{}
			 
		}
 
	longest := 0

	for _, v := range num {
		if _ ,exits := set[v-1];exits {
        continue
		}
	length :=1
	current :=v

	for  {
		if _,ok :=set[current+1];ok {
			current++
			length++
		} else {
			break
		}
	}
		if length > longest {
			longest = length
		}

	}
	fmt.Println("longest",longest)
	
 }


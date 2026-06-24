
//1. Contains Duplicate

package module2

import "fmt"

func containsDuplicate(){
	nums := []int{100,200,500,200,500}
	set := make(map[int]struct{})
	duplicate := false

	for _,v := range nums{
		 
		if _,exits := set[v];exits{
			 duplicate =true
		}
		set[v]=struct{}{}
	}
    if duplicate{
		fmt.Println("Duplicate Exits")
	}else{
		fmt.Println("Duplicate does not exits")
	}

}
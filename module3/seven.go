package module3

import (
	"fmt"
	"reflect"
)

// regular struct Declaring
type AStruct struct {
	data int
}

func StructsInterface() {
	//struct object  Creation
	obj := AStruct{
		data: 10,
	}
	fmt.Println(obj)

	//Anomous object creation
	fmt.Println(
		AStruct{
			data: 100,
		},
	)

	nestedStructFunc()
}

type Cricket struct {
	cricketRanking int
}

type FootBall struct {
	player          Cricket // Struct embedding
	footballRanking int
}

func nestedStructFunc() {
	obj1Cricket := Cricket{
		cricketRanking: 10,
	}
	obj1Football := FootBall{
		player:          obj1Cricket,
		footballRanking: 100,
	}

	fmt.Println(obj1Cricket, reflect.TypeOf(obj1Cricket))
	fmt.Println(obj1Football, reflect.TypeOf(obj1Football))
}

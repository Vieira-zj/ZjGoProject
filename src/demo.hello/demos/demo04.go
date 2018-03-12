package demos

import (
	"fmt"
)

// demo 01, init value
func init() {
	fmt.Println("start run demo04") // #2
}

func sayHello() string {
	fmt.Println("start run sayHello()") // #1
	return "hello world!"
}

// HelloMessage : test init value
var HelloMessage = sayHello()

// demo 02, struct reference
type mySubStruct struct {
	id  uint
	val string
}

type mySuperStruct struct {
	sub mySubStruct // by value
	ex  string
}

type mySuperStructRef struct {
	sub *mySubStruct // by refrence
	ex  string
}

func testStructRefValue() {
	sub := mySubStruct{
		id:  10,
		val: "ten",
	}

	super := mySuperStruct{
		sub: sub,
		ex:  "number 10",
	}
	fmt.Printf("before => sub struct: %+v\n", super)

	superRef := mySuperStructRef{
		sub: &sub,
		ex:  "number 10",
	}
	fmt.Printf("before => sub struct ref: %+v\n", superRef.sub)

	sub.val = "TEN"
	fmt.Printf("after => sub struct: %+v\n", super)
	fmt.Printf("after => sub struct Ref: %+v\n", superRef.sub)
}

// MainDemo04 : main
func MainDemo04() {
	testStructRefValue()

	fmt.Println("demo 04 done.")
}

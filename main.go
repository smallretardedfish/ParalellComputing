package main

import (
	"fmt"
)

type foo struct {
	A int
}
type bar struct {
	B int
}

func main() {
	var v interface{}
	var val foo
	v = val
	vv, _ := v.(bar)
	fmt.Println(vv.B)

}

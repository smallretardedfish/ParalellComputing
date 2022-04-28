package main

import (
	"fmt"
)

type foo struct {
	B int
}
type bar struct {
	B int
}

func main() {
	var v interface{}
	var val = foo{B: 228}
	v = val
	vv, _ := v.(bar)
	fmt.Println(vv.B)

}

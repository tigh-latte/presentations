package main

import (
	"fmt"
	"unsafe"
)

func main() {
	type person struct {
		name       string
		age        int
		registered bool
		flags      [2]byte
	}

	p := person{
		name:       "bob",
		age:        45,
		registered: true,
		flags:      [...]byte{'a', 'b'},
	}

	startAddr := unsafe.Pointer(&p)
	ageOffset := unsafe.Offsetof(p.age)

	ptr := unsafe.Add(startAddr, ageOffset)

	registered := *(*int)(ptr)
	fmt.Println(registered)

}

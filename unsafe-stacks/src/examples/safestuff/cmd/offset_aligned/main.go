package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println(unsafe.Offsetof(struct {
		registered bool    // 1 byte
		holy       bool    // 1 byte
		hell       bool    // 1 byte
		oh         bool    // 1 byte
		wow        bool    // 1 byte
		flags      [2]byte // 2 bytes
		name       string  // 16 bytes
		age        int     // 8 bytes
	}{}.name)) // Output: 8
}

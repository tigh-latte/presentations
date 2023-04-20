package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println(unsafe.Offsetof(struct {
		registered bool    // 1 byte
		name       string  // 16 bytes
		age        int     // 8 bytes
		flags      [2]byte // 2 bytes
	}{}.name)) // Output: 8
}

package main

import (
	"fmt"
)

func main() {
	fmt.Println(struct {
		registered bool    // 1 byte
		name       string  // 16 bytes
		holy       bool    // 1 byte
		age        int     // 8 bytes
		hell       bool    // 1 byte
		flags      [2]byte // 2 bytes
	}{}) // Output: 48
}

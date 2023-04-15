package main

import (
	"fmt"
)

func main() {
	fmt.Println(struct {
		registered bool    // 1 byte
		name       string  // 16 bytes
		age        int     // 8 bytes
		flags      [2]byte // 2 bytes
	}{}) // Output: 32
}

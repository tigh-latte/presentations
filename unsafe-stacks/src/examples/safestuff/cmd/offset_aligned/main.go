package main

import (
	"fmt"
)

func main() {
	fmt.Println(struct {
		name       string  // 16 bytes
		age        int     // 8 bytes
		registered bool    // 1 byte
		flags      [2]byte // 2 bytes
	}{}) // Output: 25
}

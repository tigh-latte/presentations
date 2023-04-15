package main

import (
	"fmt"
)

func main() {
	fmt.Println(struct {
		name        string  // 16 bytes
		age         int     // 8 bytes
		registered  bool    // 1 byte
		flags       [2]byte // 2 bytes
		isInitiated bool    // 1 byte
		holy        bool    // 1 byte
		hell        bool    // 1 byte
		oh          bool    // 1 byte
		wow         bool    // 1 byte
	}{}) // Output: 32
}

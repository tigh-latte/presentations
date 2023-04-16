package main

import (
	"fmt"
	"strings"
	"unsafe"
)

type MyCoolString string

func main() {
	safe()
	unsafeCast()
}

func safe() {
	v := []MyCoolString{"hello", "there", "oh", "wow"}

	ss := make([]string, len(v))

	for i := 0; i < len(v); i++ {
		ss[i] = string(v[i])
	}

	fmt.Println(strings.Join(ss, ","))
}

func unsafeCast() {
	v := []MyCoolString{"hello", "there", "oh", "wow"}

	ss := *(*[]string)(unsafe.Pointer(&v))
	fmt.Println(strings.Join(ss, ","))
}

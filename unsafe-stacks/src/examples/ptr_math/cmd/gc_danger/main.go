package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

func main() {
	z := func() *int {
		return new(int)
	}()

	fmt.Println("mem z:", z, *z)

	addrZ := uintptr(unsafe.Pointer(z))

	addrZ += 2
	addrZ -= 2

	runtime.GC()

	var (
		p = 7
		q = 8
		s = 9
	)

	*(*int)(unsafe.Pointer(addrZ)) = 3

	fmt.Println("mem p:", &p, p)
	fmt.Println("mem q:", &q, q)
	fmt.Println("mem s:", &s, s)
}

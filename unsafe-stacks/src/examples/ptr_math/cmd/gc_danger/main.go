package main

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"unsafe"
)

func main() {
	log.SetOutput(io.Discard)

	z := func() *int {
		return new(int)
	}()

	fmt.Println("mem z:", z, *z)

	ptr := uintptr(unsafe.Pointer(z))

	runtime.GC()

	p, q, s := 7, 8, 9

	*(*int)(unsafe.Pointer(ptr)) = -19

	log.Print((*[8]byte)(unsafe.Pointer(&q)))
	log.Print((*[8]byte)(unsafe.Pointer(&ptr)))

	fmt.Println("mem p:", &p, p)
	fmt.Println("mem q:", &q, q)
	fmt.Println("mem s:", &s, s)
}

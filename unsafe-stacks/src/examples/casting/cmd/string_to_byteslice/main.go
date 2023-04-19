package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
}

func safe() {
	s := "Hello world"
	bb := []byte(s)
	fmt.Println(bb)
}

func unsafeCast() {
	s := "Hello world"
	bb := *(*[]byte)(unsafe.Pointer(&s))
	fmt.Println(bb) // Output: [72 101 108 108 111 32 119 111 114 108 100]
}

func unsafeSlice() {
	s := "Hello world"
	bb := unsafe.Slice(unsafe.StringData(s), len(s))
	fmt.Println(bb) // Output: [72 101 108 108 111 32 119 111 114 108 100]
}

func unsafeHeader() {
	s := "Hello world"

	var bb []byte
	bbHdr := (*reflect.SliceHeader)(unsafe.Pointer(&bb))
	sHdr := (*reflect.StringHeader)(unsafe.Pointer(&s))

	bbHdr.Data = sHdr.Data
	bbHdr.Len = sHdr.Len
	bbHdr.Cap = sHdr.Len

	fmt.Println(bb) // Output: [72 101 108 108 111 32 119 111 114 108 100]
}

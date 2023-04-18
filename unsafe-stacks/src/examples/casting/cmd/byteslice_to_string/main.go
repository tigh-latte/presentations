package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
}

func safe() {
	bb := []byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}
	s := string(bb)
	fmt.Println(s)
}

func unsafeCast() {
	bb := []byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}
	s := (*string)(unsafe.Pointer(&bb))
	fmt.Println(*s) // Output: Hello world!
}

func unsafeString() {
	bb := []byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}
	s := unsafe.String(unsafe.SliceData(bb), len(bb))
	fmt.Println(s) // Output: Hello world!
}

func unsafeHeader() {
	bb := []byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}

	var s string
	sHdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bbHdr := (*reflect.SliceHeader)(unsafe.Pointer(&bb))

	sHdr.Data = bbHdr.Data
	sHdr.Len = bbHdr.Len

	fmt.Println(s) // Output: Hello world!
}

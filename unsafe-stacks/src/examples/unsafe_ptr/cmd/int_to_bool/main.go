package main

import (
	"fmt"
	"unsafe"
)

func main() {
	i := new(int64)
	b := (*bool)(unsafe.Pointer(i))
	f := (*float64)(unsafe.Pointer(i))
	a := (*[8]byte)(unsafe.Pointer(i))

	fmt.Println(*b) // Output: false

	*i = 10
	fmt.Println(*b) // Output: true
	fmt.Println(*f) // Output: 5e-323
	fmt.Println(*a) // Output: [10 0 0 0 0 0 0 0]

	*i = 256
	fmt.Println(*b) // Output: false
	fmt.Println(*f) // Output: 1.265e-321
	fmt.Println(*a) // Output: [0 1 0 0 0 0 0 0]

	*i = -10
	fmt.Println(*b) // Output: true
	fmt.Println(*f) // Output: NaN
	fmt.Println(*a) // Output: [246 255 255 255 255 255 255 255]
}

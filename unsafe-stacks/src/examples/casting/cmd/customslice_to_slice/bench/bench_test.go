package bench

import (
	"fmt"
	"testing"
	"unsafe"
)

var global []string

func Benchmark_SafeSliceCast(b *testing.B) {
	type testType string

	for _, test := range []int{10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("%d elems", test), func(b *testing.B) {
			ww := make([]testType, test)
			for i := range ww {
				ww[i] = testType("hello")
			}

			b.ResetTimer()

			var out []string
			for i := 0; i < b.N; i++ {
				ss := make([]string, len(ww))
				for i := range ww {
					ss[i] = string(ww[i])
				}

				out = ss
			}

			global = out
		})
	}
}

func Benchmark_UnsafeSliceCast(b *testing.B) {
	type testType string

	for _, test := range []int{10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("%d elems", test), func(b *testing.B) {
			ww := make([]testType, test)
			for i := range ww {
				ww[i] = testType("hello")
			}

			b.ResetTimer()
			var out []string
			for i := 0; i < b.N; i++ {
				ss := *(*[]string)(unsafe.Pointer(&ww))
				out = ss
			}

			global = out
		})
	}
}

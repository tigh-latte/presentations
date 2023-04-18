package bench_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

var global string

func Benchmark_SafeCast(b *testing.B) {
	for _, test := range []int{1, 10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("%d len", test), func(b *testing.B) {
			bb := bytes.Repeat([]byte{'a'}, test)

			var out string

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s := string(bb)
				out = s
			}

			global = out
		})
	}
}

func Benchmark_UnsafeCast(b *testing.B) {
	for _, test := range []int{1, 10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("%d len", test), func(b *testing.B) {
			bb := bytes.Repeat([]byte{'a'}, test)

			var out string

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s := *(*string)(unsafe.Pointer(&bb))
				out = s
			}

			global = out
		})
	}
}

func Benchmark_UnsafeString(b *testing.B) {
	for _, test := range []int{1, 10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("%d len", test), func(b *testing.B) {
			bb := bytes.Repeat([]byte{'a'}, test)

			var out string

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s := unsafe.String(unsafe.SliceData(bb), len(bb))
				out = s
			}

			global = out
		})
	}
}

func Benchmark_UnsafeHeader(b *testing.B) {
	for _, test := range []int{1, 10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("%d len", test), func(b *testing.B) {
			bb := bytes.Repeat([]byte{'a'}, test)

			var out string

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var s string

				sHdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
				bbHdr := (*reflect.SliceHeader)(unsafe.Pointer(&bb))
				sHdr.Data = bbHdr.Data
				sHdr.Len = bbHdr.Len

				out = s
			}

			global = out
		})
	}
}

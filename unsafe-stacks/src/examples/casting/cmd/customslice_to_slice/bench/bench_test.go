package bench

import (
	"testing"
	"unsafe"
)

var global []string

func Benchmark_SafeSliceCast(b *testing.B) {
	type testType string

	tests := []struct {
		name  string
		total int
	}{{
		name:  "10 elems",
		total: 10,
	}, {
		name:  "100 elems",
		total: 100,
	}, {
		name:  "1000 elems",
		total: 1000,
	}, {
		name:  "10000 elems",
		total: 10000,
	}, {
		name:  "100000 elems",
		total: 100000,
	}}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			ww := make([]testType, test.total)
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

	tests := []struct {
		name  string
		total int
	}{{
		name:  "10 elems",
		total: 10,
	}, {
		name:  "100 elems",
		total: 100,
	}, {
		name:  "1000 elems",
		total: 1000,
	}, {
		name:  "10000 elems",
		total: 10000,
	}, {
		name:  "100000 elems",
		total: 100000,
	}}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			ww := make([]testType, test.total)
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

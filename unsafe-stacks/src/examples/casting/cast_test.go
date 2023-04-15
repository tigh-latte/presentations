package cast_test

import (
	"bytes"
	"testing"
	"unsafe"
)

var global string

func Benchmark_Safe1(b *testing.B) {
	bb := []byte{'a'}

	var w string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w = string(bb)
	}

	global = w
}

func Benchmark_Safe100(b *testing.B) {
	bb := bytes.Repeat([]byte{'a'}, 100)

	var w string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w = string(bb)
	}

	global = w
}

func Benchmark_Safe10000(b *testing.B) {
	bb := bytes.Repeat([]byte{'a'}, 10000)

	var w string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w = string(bb)
	}

	global = w
}

func Benchmark_Safe100000(b *testing.B) {
	bb := bytes.Repeat([]byte{'a'}, 100000)

	var w string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w = string(bb)
	}

	global = w
}

func Benchmark_Unsafe1(b *testing.B) {
	bb := []byte{'a'}

	var w string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w = unsafe.String(unsafe.SliceData(bb), len(bb))
	}

	global = w
}

func Benchmark_Unsafe100(b *testing.B) {
	bb := bytes.Repeat([]byte{'a'}, 100)

	var w string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w = unsafe.String(unsafe.SliceData(bb), len(bb))
	}

	global = w
}

func Benchmark_Unsafe10000(b *testing.B) {
	bb := bytes.Repeat([]byte{'a'}, 10000)

	var w string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w = unsafe.String(unsafe.SliceData(bb), len(bb))
	}

	global = w
}

func Benchmark_Unsafe100000(b *testing.B) {
	bb := bytes.Repeat([]byte{'a'}, 100000)

	var w string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w = unsafe.String(unsafe.SliceData(bb), len(bb))
	}

	global = w
}

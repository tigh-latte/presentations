package util_test

import (
	"strings"
	"testing"

	"github.com/tigh-latte/strings/util"
)

var catch string

func Benchmark_CloneSafe(b *testing.B) {
	s := strings.Repeat(".", 42)
	for i := 0; i < b.N; i++ {
		catch = util.CloneSafe(s)
	}
}

func Benchmark_CloneUnsafe(b *testing.B) {
	s := strings.Repeat(".", 42)
	for i := 0; i < b.N; i++ {
		catch = util.CloneUnsafe(s)
	}
}

package crypt_test

import (
	"crypto/aes"
	"testing"

	"github.com/tigh-latte/heap/crypt"
)

func Benchmark_EncryptSafe(b *testing.B) {
	var key [16]byte
	cipher, _ := aes.NewCipher(key[:])
	for i := 0; i < b.N; i++ {
		crypt.EncryptSafe(cipher, "hello")
	}
}

func Benchmark_EncryptUnsafe(b *testing.B) {
	var key [16]byte
	cipher, _ := aes.NewCipher(key[:])
	for i := 0; i < b.N; i++ {
		crypt.EncryptUnsafe(cipher, "hello")
	}
}

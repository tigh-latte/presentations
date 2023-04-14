package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	_ "unsafe"
)

func EncryptSafe(block cipher.Block, message string) (string, error) {
	byteMsg := []byte(message)

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("could not encrypt: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func EncryptUnsafe(block cipher.Block, message string) (string, error) {
	b := make([]byte, len(message))
	copy(b, message)

	cipherText := make([]byte, aes.BlockSize+len(b))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("could not encrypt: %v", err)
	}

	stream := newCFBEncrypter(block, iv)
	xorKeyStream(stream, cipherText[aes.BlockSize:], b)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

//go:noescape
//go:linkname newCFBEncrypter github.com/tigh-latte/heap/crypt.__newCFBEncrypter
func newCFBEncrypter(b cipher.Block, iv []byte) cipher.Stream
func __newCFBEncrypter(b cipher.Block, iv []byte) cipher.Stream {
	return cipher.NewCFBEncrypter(b, iv)
}

//go:noescape
//go:linkname xorKeyStream github.com/tigh-latte/heap/crypt.__xorKeyStream
func xorKeyStream(s cipher.Stream, dst, src []byte)
func __xorKeyStream(s cipher.Stream, dst, src []byte) {
	s.XORKeyStream(dst, src)
}

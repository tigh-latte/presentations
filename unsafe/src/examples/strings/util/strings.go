package util

import "unsafe"

func CloneSafe(s string) string {
	if len(s) == 0 {
		return ""
	}

	b := make([]byte, len(s))
	copy(b, s)

	return string(b)
}

func CloneUnsafe(s string) string {
	if len(s) == 0 {
		return ""
	}

	b := make([]byte, len(s))
	copy(b, s)

	return unsafe.String(&b[0], len(b))
}

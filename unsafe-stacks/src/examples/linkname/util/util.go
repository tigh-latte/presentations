package util

import (
	"math/rand"
)

func PositiveRand(i int) int {
	o := rand.Int() % i
	if o < 0 {
		return -o
	}
	return o
}

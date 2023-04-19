package util_test

import (
	"testing"

	_ "unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/tigh-latte/linkname/util"
)

//go:linkname randInt math/rand.Int
func randInt() int {
	return 9
}

func Test_PositiveRan(t *testing.T) {
	assert.Equal(t, 2, util.PositiveRand(7))
}

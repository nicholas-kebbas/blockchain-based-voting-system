package tests

import (
	"fmt"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p3"
	"testing"
)

func TestRandom(t *testing.T) {
	p3.RandSeed()
	fmt.Println(p3.GenerateRandomString(16))
}
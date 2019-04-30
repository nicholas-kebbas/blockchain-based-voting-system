package tests

import (
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p3"
	"testing"
)


func TestNonces(t *testing.T) {
	p3.RandSeed()
	p3.StartTryingNonces()
}
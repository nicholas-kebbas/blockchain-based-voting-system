package tests

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestSignature(t *testing.T) {

	/* Create and Sign */

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature_p: (0x%x, 0x%x)\n", r, s)
	/* Validate Truth */
	valid := ecdsa.Verify(&privateKey.PublicKey, hash[:], r, s)
	fmt.Println("signature_p verified:", valid)
	if !valid {
		t.Fail()
	}
	//testByteArray := []byte("test")
	//signature_p.CreateSignature(testByteArray)



}
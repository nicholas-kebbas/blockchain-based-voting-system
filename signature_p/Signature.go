package signature_p

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"errors"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"math/big"
	"os"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p2"
)

var PUBLIC_KEY = []byte{}
var PRIVATE_KEY = new(ecdsa.PrivateKey)
var SIGNATURE = []byte{}

/* Need to make a ECDSASignature struct to verify */
type ECDSASignature struct {
	R, S *big.Int
}

/* Need to implement basic signatures first. */

/* In order to sign data, we need the data to sign.

So we need to take that as input. This is probably MPT

*/

/* Take the hash of the MPT and encrypt it */
/* We will need to add a Signature and a Public Key parameter to the Block */

func SignTransaction(value string) {
	transaction := []byte (value)
	r := big.NewInt(0)
	s := big.NewInt(0)
	serr := errors.New("Error")
	/* Returns Big Ints r and s
	*/
	r, s, serr = ecdsa.Sign(crand.Reader, PRIVATE_KEY, transaction)
	if serr != nil {
		fmt.Println("Error")
		os.Exit(1)
	}

	/* Need to figure out how to get the r and s values from the signature_p on the blockchain */

	/*
	The signature_p is a combination of the author's private key and the content of
	the document it certifies
	 */
	SIGNATURE = r.Bytes()
	SIGNATURE = append(SIGNATURE, s.Bytes()...)
}

/* To verify the data, we need
1. The data that was signed
2. The signature_p
3. The Public Key

So we need to add that to the Block
Will likely take as input the signature_p
*/
func VerifySignature(block p2.Block) bool{
	e := &ECDSASignature{}
	_, err := asn1.Unmarshal([]byte(block.Header.Signature), e)
	if err != nil {
		fmt.Println("Error Unmarshaling Block")
		return false
	}
	/* Currently a bug here */
	//verified := ecdsa.Verify(&block.Header.PublicKey, []byte(block.GetMptRoot()), e.R, e.S)
	return true
}

/* Generate the keys we need for the addresses, signature_p, etc.

In elliptic curve based algorithms, public keys are points on a curve.
A public key is a combination of X, Y coordinates.

PUBLIC_KEY should be stored as a byte array since it's easier to
work with that way.

*/


func GeneratePublicAndPrivateKey() {
	c := elliptic.P256()
	PRIVATE_KEY, _ = ecdsa.GenerateKey(c, crand.Reader)
	PUBLIC_KEY = append(PRIVATE_KEY.PublicKey.X.Bytes(), PRIVATE_KEY.PublicKey.Y.Bytes()...)
}

/* Hash the public key to display on the blockchain. This is how BTC does it */
/* TODO: Implement this function for additional ***security*** and ***privacy***! */
func HashPublicKey(publicKey []byte) []byte {
	publicSHA256 := sha256.Sum256(publicKey)
	ripemd160Hasher := ripemd160.New()
	_, err := ripemd160Hasher.Write(publicSHA256[:])

	if err != nil {
		fmt.Println("Cannot Hash Public Key")
		os.Exit(1)
	}

	hashedPublicKey := ripemd160Hasher.Sum(nil)
	return hashedPublicKey
}

/* Merges transaction with other transactions on the chain to maintain anonymity */
func RingSignature() string {
	ringSignature := ""

	return ringSignature

}
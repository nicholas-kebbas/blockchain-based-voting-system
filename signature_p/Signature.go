package signature_p

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p2"
	"golang.org/x/crypto/ripemd160"
	"log"
	"math/big"
	"os"
)


var PUBLIC_KEY= []byte{}
var PUBLIC_KEY_HEX = ""
var PRIVATE_KEY = new(ecdsa.PrivateKey)
var PRIVATE_KEY_BYTES = []byte{}
var PRIVATE_KEY_HEX = ""
var SIGNATURE = []byte{}
var SIGNATURE_HEX = ""

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

/* Return the hex signature */
func SignTransaction(value string) []byte{
	transaction := []byte (value)
	/* Hash this again. May only need to hash it once but lib wants it as 32 bits */
	hash := crypto.Keccak256Hash(transaction)
	/* Generate Signature as well */
	signature, err := crypto.Sign(hash.Bytes(), PRIVATE_KEY)
	if err != nil {
		log.Fatal(err)
	}
	SIGNATURE = signature
	SIGNATURE_HEX = hexutil.Encode(signature)
	fmt.Println("Signature hex")
	fmt.Println(SIGNATURE_HEX)
	/* Try the hex version but shouldn't be any different */
	return signature

}

/* To verify the data, we need
1. The data that was signed
2. The signature_p
3. The Public Key

So we need to add that to the Block
Will likely take as input the signature_p
*/
func VerifySignature(block p2.Block) bool {
	fmt.Println("Verify: Signature in Block")
	fmt.Println(block.Header.Signature)
	decodedSig, _ := hexutil.Decode(block.Header.Signature)
	hash := crypto.Keccak256Hash([]byte(block.GetMptRoot()))
	sigPublicKey, _ := crypto.Ecrecover(hash.Bytes(), decodedSig)
	signatureNoRecoverID := decodedSig[:len(decodedSig)-1] // remove recovery ID
	verified := crypto.VerifySignature(sigPublicKey, hash.Bytes(), []byte(signatureNoRecoverID))

	fmt.Println("Verified Signature")
	fmt.Println(verified)
	return verified
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
	PRIVATE_KEY_BYTES = PRIVATE_KEY.D.Bytes()
	PRIVATE_KEY_HEX = hexutil.Encode(PRIVATE_KEY_BYTES)[2:]
	publicKey := PRIVATE_KEY.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	PUBLIC_KEY = crypto.FromECDSAPub(publicKeyECDSA)
	PUBLIC_KEY_HEX = hexutil.Encode(PUBLIC_KEY)[4:]
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
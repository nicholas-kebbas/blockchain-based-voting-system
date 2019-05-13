package tests

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p1"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p3/data"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/signature_p"
	"testing"
)

func TestSignature(t *testing.T) {
	/* Create block and Sign */
	newBlockChain := data.NewBlockChain()
	signature_p.GeneratePublicAndPrivateKey()
	mpt := p1.MerklePatriciaTrie{}
	mpt.Initial()
	/* First block does not need to be verified, rest do */
	mpt.Insert("1", "Origin")
	hexPubKey := hexutil.Encode(signature_p.PUBLIC_KEY)
	newBlockChain.GenBlock(mpt, hexPubKey)
	/* So now we have the signed thing, let's verify */
	block, _ := newBlockChain.Get(1)
	match := signature_p.VerifySignature(block[0])
	fmt.Println("Match")
	fmt.Println(match)
}

func TestSignatureBasic(t *testing.T) {
	privateKey, _ := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")

	data := []byte("hello")
	hash1 := crypto.Keccak256Hash(data)
	fmt.Println(hash1.Hex()) // 0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8

	signature, _ := crypto.Sign(hash1.Bytes(), privateKey)
	sigPublicKey, _ := crypto.Ecrecover(hash1.Bytes(), signature)
	/* Run the same hash on same value */
	hash2 := crypto.Keccak256Hash([]byte("hello"))
	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery ID
	verified := crypto.VerifySignature(sigPublicKey, hash2.Bytes(), signatureNoRecoverID)
	fmt.Println(verified)
}
package data

import (
	"errors"
	"fmt"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p1"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p2"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/signature_p"
	"sync"
)

type SyncBlockChain struct {
	bc p2.BlockChain
	mux sync.Mutex
}

func NewBlockChain() SyncBlockChain {
	return SyncBlockChain{bc: p2.NewBlockChain()}
}

/* Returns false if not found */
func(sbc *SyncBlockChain) Get(height int32) ([]p2.Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	if sbc.bc.Get(height) != nil {
		return sbc.bc.Get(height), true
	}
	return sbc.bc.Get(height), false
}

func(sbc *SyncBlockChain) GetLength() (height int32) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Length
}

func(sbc *SyncBlockChain) GetBlock(height int32, hash string) (p2.Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	blocks := sbc.bc.Get(height)
	i := 0
	for i = 0; i < len(blocks); i++ {
		if blocks[i].Header.Hash == hash {
			return blocks[i], true
		}
	}
	emptyBlock := p2.Block{}
	return emptyBlock, false
}

func(sbc *SyncBlockChain) Insert(block p2.Block) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	sbc.bc.Insert(block)
}

/* Check if this block is found in the chain */
func(sbc *SyncBlockChain) CheckParentHash(insertBlock p2.Block) bool {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	if sbc.bc.CheckForHash(insertBlock.Header.ParentHash) {
		return true
	}
	return false
}

func(sbc *SyncBlockChain) UpdateEntireBlockChain(blockChainJson string) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	p2.DecodeFromJSON(&sbc.bc, blockChainJson)
}

func(sbc *SyncBlockChain) BlockChainToJson() (string, error) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	jsonBlockChain := sbc.bc.EncodeToJSON()
	return jsonBlockChain,  errors.New("There was an error")
}

/* Create a new block to add to Sync. Blockchain */
func(sbc *SyncBlockChain) GenBlock(mpt p1.MerklePatriciaTrie, public_key []byte, signature []byte) p2.Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	height := sbc.bc.Length
	if height == 0 {
		sbc.bc = p2.NewBlockChain()
		/* We have our Public Key Stored in signature_p */

		signature = signature_p.SignTransaction(mpt.GetRoot())
		fmt.Println("Signature in Gen Block")
		fmt.Println(signature)
		newBlock := p2.Initial(1, 123, "Genesis", mpt, public_key, signature)

		sbc.bc.Insert(newBlock)
		return newBlock
	}
	latestBlock := sbc.bc.Get(height)[0]

	/* Add height because it's a new block, child block of parent */
	signature = signature_p.SignTransaction(mpt.GetRoot())
	fmt.Println("Signature in Gen Block")
	fmt.Println(signature)
	newBlock := p2.Initial(height + 1, 123, latestBlock.Header.Hash, mpt, public_key, signature)
	sbc.bc.Insert(newBlock)
	return newBlock
}

func (sbc *SyncBlockChain) GetLatestBlocks() []p2.Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Get(sbc.bc.Length)
}

func (sbc *SyncBlockChain) GetParentBlock(block p2.Block) (p2.Block, error) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetParentBlock(block)
}

func(sbc *SyncBlockChain) Show() string {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Show()
}

func(sbc *SyncBlockChain) ShowMPT() string {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.ShowMPT()
}
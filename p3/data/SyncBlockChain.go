package data

import (
	"errors"
	"fmt"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p1"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p2"
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
	sbc.bc.Insert(block)
	sbc.mux.Unlock()
}

/* Check if this block is found in the chain */
func(sbc *SyncBlockChain) CheckParentHash(insertBlock p2.Block) bool {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	fmt.Println("Entering sbc.bc.CheckParentHash, checking for")
	fmt.Println(insertBlock.Header.ParentHash)
	if sbc.bc.CheckForHash(insertBlock.Header.ParentHash) {
		fmt.Println("Parent Hash is found in CheckParentHash(). Returning True.")
		return true
	}
	fmt.Println("Parent Hash is not found in CheckParentHash(). Returning False.")
	return false
}

func(sbc *SyncBlockChain) UpdateEntireBlockChain(blockChainJson string) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	fmt.Println("Update Entire Block chain:")
	fmt.Println(blockChainJson)
	p2.DecodeFromJSON(&sbc.bc, blockChainJson)
	fmt.Println("Update Entire Block chain Blockchain after Decode:")
	fmt.Println(sbc.bc)
}

/* TODO: Fix error handling on this one */
func(sbc *SyncBlockChain) BlockChainToJson() (string, error) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	jsonBlockChain := sbc.bc.EncodeToJSON()
	return jsonBlockChain,  errors.New("There was an error")
}

/* TODO: We're iterating Length here because that needs to be done for this implementation,
but this should probably be done in the Blockchain class. */
func(sbc *SyncBlockChain) GenBlock(mpt p1.MerklePatriciaTrie) p2.Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	height := sbc.bc.Length
	if height == 0 {
		fmt.Println("HEIGHT IS 0.")
		fmt.Println("This should only show up in first node")
		sbc.bc = p2.NewBlockChain()
		newBlock := p2.Initial(1, 123, "Genesis", mpt)
		sbc.bc.Insert(newBlock)
		return newBlock
	}
	/* -1 because we're getting the hash of the parent. But this should just be height, maybe*/
	fmt.Println("height in GENERATE BLOCK")
	fmt.Println(height)
	latestBlock := sbc.bc.Get(height)[0]
	newBlock := p2.Initial(height + 1, 123, latestBlock.Header.Hash, mpt)
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
package data

import (
	"errors"
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

func(sbc *SyncBlockChain) GetBlock(height int32, hash string) (p2.Block, bool) {
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
	if sbc.bc.CheckForHash(insertBlock.Header.ParentHash) {
		return true
	}
	return false
}

func(sbc *SyncBlockChain) UpdateEntireBlockChain(blockChainJson string) {
	p2.DecodeFromJSON(&sbc.bc, blockChainJson)
}

/* TODO: Fix error handling on this one */
func(sbc *SyncBlockChain) BlockChainToJson() (string, error) {
	jsonBlockChain := sbc.bc.EncodeToJSON()
	return jsonBlockChain,  errors.New("need to fix this error handling")
}

/* TODO: Generate Block only having the MPT. */
func(sbc *SyncBlockChain) GenBlock(mpt p1.MerklePatriciaTrie) p2.Block {
	sbc.mux.Lock()
	height := sbc.bc.Length
	if height == 0 {
		sbc.bc = p2.NewBlockChain()
		newBlock := p2.Initial(height, 123, "Genesis", mpt)
		sbc.bc.Insert(newBlock)
		sbc.mux.Unlock()
		return newBlock
	}
	latestBlock := sbc.bc.Get(height)[0]
	newBlock := p2.Initial(height, 123, latestBlock.Header.Hash, mpt)
	sbc.bc.Insert(newBlock)
	sbc.mux.Unlock()
	return newBlock
}

func(sbc *SyncBlockChain) Show() string {
	return sbc.bc.Show()
}
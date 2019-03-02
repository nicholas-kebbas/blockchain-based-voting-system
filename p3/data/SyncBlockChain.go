package data

import (
	"../../p1"
	"../../p2"
	"sync"
)

type SyncBlockChain struct {
	bc p2.BlockChain
	mux sync.Mutex
}

func NewBlockChain() SyncBlockChain {
	return SyncBlockChain{bc: p2.NewBlockChain()}
}

func(sbc *SyncBlockChain) Get(height int32) ([]p2.Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Get(height)
}

func(sbc *SyncBlockChain) GetBlock(height int32, hash string) (p2.Block, bool) {}

func(sbc *SyncBlockChain) Insert(block p2.Block) {
	sbc.mux.Lock()
	sbc.bc.Insert(block)
	sbc.mux.Unlock()
}

func(sbc *SyncBlockChain) CheckParentHash(insertBlock p2.Block) bool {}

func(sbc *SyncBlockChain) UpdateEntireBlockChain(blockChainJson string) {}

func(sbc *SyncBlockChain) BlockChainToJson() (string, error) {}

func(sbc *SyncBlockChain) GenBlock(mpt p1.MerklePatriciaTrie) p2.Block {}

func(sbc *SyncBlockChain) Show() string {
	return sbc.bc.Show()
}
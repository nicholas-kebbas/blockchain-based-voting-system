package p2

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"sort"
)

/* The value is a list so it can handle the forks */
type BlockChain struct {
	Length int32
	Chain map[int32][]Block
}

type JsonBlockChain struct {
	JsonBlockArray []string
}

/**
Create new Blockchain object.
 */
func NewBlockChain() BlockChain {
	blockchain := BlockChain{}
	blockchain.Length = 0
	blockchain.Chain = make(map[int32][]Block)
	return blockchain
}

/**
Get array of blocks at given height. If empty, return nil.
 */
func (blockchain *BlockChain) Get(height int32) []Block {
	blockArray := blockchain.Chain[height]
	if len(blockArray) > 0 {
		return blockArray
	}
	return nil
}

/* Check if hash exists in blockchain */
func (blockchain *BlockChain) CheckForHash(hash string) bool {
	var i int32 = 0
	for i = 0; i < blockchain.Length; i++ {
		for z := 0; z < len(blockchain.Get(i)); z++ {
			fmt.Println("Check for HASH header hash")
			fmt.Println(blockchain.Get(i)[z].Header.Hash)
			if blockchain.Get(i)[z].Header.Hash == hash {
				fmt.Println("Found the hash")
				return true
			}
		}
	}
	fmt.Println("Can't find the parent in Blockchain CheckforHash()")
	return false
}


/* Insert into the chain. mapping the height to the block itself */
/*  If the list has already contained that block's hash, ignore it because we don't store duplicate blocks;
if not, insert the block into the list.
 */
func (blockchain *BlockChain) Insert(block Block) {
	/* Check for block's hash in chain. Need to start at end and work way back */
	for i := blockchain.Length; i > 0; i-- {
		blockIterator := blockchain.Chain[i]
		for _, b := range blockIterator {
			if block.Header.Hash == b.Header.Hash {
				fmt.Println("Duplicate Block")
				return
			}
		}
	}
	/* else it's not in the chain, so add it and adjust BC Length */

	/* Need to get height of block being entered */
	blockArray := blockchain.Chain[block.Header.Height]
	blockArray = append(blockArray, block)

	/* Update length if new largest height */
	if block.Header.Height > blockchain.Length {
		blockchain.Length = block.Header.Height
	}
	blockchain.Chain[block.Header.Height] = blockArray

}

/**
Function of the blockchain object. Turns itself into json representation. Returns String.
 */
func (blockchain *BlockChain) EncodeToJSON() string {
	jsonBlockArray := []JsonBlock{}
	i := int32(0)
	for i = 0; i <= blockchain.Length; i++ {
		blockIterator := blockchain.Chain[i]
		for z := range blockIterator {
			block := blockIterator[z]
			blockString := block.EncodeToJSON()
			jsonBytes := []byte(blockString)
			jsonBlock2 := JsonBlock{}
			err := json.Unmarshal(jsonBytes, &jsonBlock2)
			if err != nil {
				fmt.Println("Error")
			}
			jsonBlockArray = append(jsonBlockArray, jsonBlock2)
		}
	}


	bs2, _:= json.Marshal(jsonBlockArray)
		return string(bs2)
	}

/**
Takes as input the blockchain that will have blocks added to it, and a json representation of blocks.
This might be causing a problem.
 */
func DecodeFromJSON(blockchain *BlockChain, jsonString string) {
	b := Block{}
	var jsonEncodedBlocks []JsonBlock
	err := json.Unmarshal([]byte(jsonString), &jsonEncodedBlocks)
	if err != nil {
		fmt.Println("Error")
	}
	/** Loop through and turn jsonBlocks into blocks, then add blocks to chain **/
	for i := range jsonEncodedBlocks {
	jsonString, _ := json.Marshal(jsonEncodedBlocks[i])
		b = b.DecodeFromJson(string(jsonString))
		blockchain.Insert(b)
	}
	fmt.Println("Decode From: ")
	fmt.Println(blockchain)
}

func (jbc JsonBlockChain) String() string {
	return fmt.Sprintf("Blockchain=%v", jbc)
}

func (bc *BlockChain) Show() string {
	rs := ""
	var idList []int
	for id := range bc.Chain {
		idList = append(idList, int(id))
	}
	sort.Ints(idList)
	for _, id := range idList {
		var hashs []string
		for _, block := range bc.Chain[int32(id)] {
			hashs = append(hashs, block.Header.Hash + "<=" + block.Header.ParentHash)
		}
		sort.Strings(hashs)
		rs += fmt.Sprintf("%v: ", id)
		for _, h := range hashs {
			rs += fmt.Sprintf("%s, ", h)
		}
		rs += "\n"
	}
	sum := sha3.Sum256([]byte(rs))
	rs = fmt.Sprintf("This is the BlockChain: %s\n", hex.EncodeToString(sum[:])) + rs
	return rs
}



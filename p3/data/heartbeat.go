package data

import (
	"encoding/json"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p1"
)

/* Heartbeat is the JSON representation of the data we need to send to the other blockchains
 */
type HeartBeatData struct {
	IfNewBlock  bool   `json:"ifNewBlock"`
	Id          int32  `json:"id"`
	BlockJson   string `json:"blockJson"`
	PeerMapJson string `json:"peerMapJson"`
	Addr        string `json:"addr"`
	Hops        int32  `json:"hops"`
}

func NewHeartBeatData(ifNewBlock bool, id int32, blockJson string, peerMapJson string, addr string) HeartBeatData {
	heartBeatData := HeartBeatData{ifNewBlock,
		id,
		blockJson,
		peerMapJson,
		addr,
		2}
	return heartBeatData
}

/* Create a new instance of HeartBeatData, then decide whether to create a new block and send it to other peers.
*/
func PrepareHeartBeatData(sbc *SyncBlockChain, selfId int32, peerMapJson string, addr string, verified bool, trie p1.MerklePatriciaTrie) HeartBeatData {
	/* Randomly decide whether to create new block and send to peers. */
	if verified == true {
		heartBeatData := NewHeartBeatData(true, selfId, " ", peerMapJson, addr)
		/* Just get the first one in that array for now since we don't know what to do w/ forks */
		/* This is adding to own fine. Maybe overwriting parent */
		newBlock := sbc.GenBlock(trie)
		heartBeatData.BlockJson = newBlock.EncodeToJSON()
		return heartBeatData
	} else {
		heartBeatData := NewHeartBeatData(false, selfId, " ", peerMapJson, addr)
		return heartBeatData
	}
}

func (heartbeat *HeartBeatData) HeartBeatToJson() string {
	jsonHeartBeat, _ := json.Marshal(heartbeat)
	return string(jsonHeartBeat)
}
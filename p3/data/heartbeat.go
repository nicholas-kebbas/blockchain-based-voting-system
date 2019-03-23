package data

import (
	"math/rand"
)

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
		3}
	return heartBeatData

}

/* Create a new instance of HeartBeatData, then decide whether to create a new block and send it to other peers.
  These arguments are currently wrong. Not sure what to make them.
*/
func PrepareHeartBeatData(sbc *SyncBlockChain, selfId int32, peerMapJson string, addr string) HeartBeatData {
	heartBeatData := NewHeartBeatData(true, selfId, " ", peerMapJson, addr);

	/* Randomly decide whether to send it to new peers. If true create block. */
	if rand2() == true {
		return heartBeatData
	}

	return heartBeatData
}

func rand2() bool {
	return rand.Int31()&0x01 == 0
}
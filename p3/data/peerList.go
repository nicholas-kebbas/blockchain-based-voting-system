package data

import (
	"encoding/json"
	// "encoding/json"
	"fmt"
	"reflect"
	// "sort"
	// "strings"
	"sync"
)

type PeerList struct {
	selfId int32
	peerMap map[string]int32
	maxLength int32
	mux sync.Mutex
}

func NewPeerList(id int32, maxLength int32) PeerList {
	/* Create new peer map */
	peerList := PeerList{}
	peerList.selfId = id
	peerList.peerMap = make(map[string]int32)
	peerList.maxLength = maxLength
	return peerList
}

func(peers *PeerList) Add(addr string, id int32) {
	peers.peerMap[addr] = id
}

func(peers *PeerList) Delete(addr string) {
	delete(peers.peerMap, addr)
}

func(peers *PeerList) Rebalance() {

}

func(peers *PeerList) Show() string {
	show := "Show"
	fmt.Println(show)
	return show
}

func(peers *PeerList) Register(id int32) {
	peers.selfId = id
	fmt.Printf("SelfId=%v\n", id)
}

/* Return copy of peer list presumably */
func(peers *PeerList) Copy() map[string]int32 {
	newMap := make(map[string]int32)
	for k,v := range peers.peerMap {
		newMap[k] = v
	}
	return newMap
}

func(peers *PeerList) GetSelfId() int32 {
	return peers.selfId
}

/* TODO: Fix error checking */
func(peers *PeerList) PeerMapToJson() (string, error) {
	jsonPeerMap, err := json.Marshal(peers.peerMap)
	return string(jsonPeerMap), err
}

/* Looks like this will take peerMap as json String and add it to existing peer map */
func(peers *PeerList) InjectPeerMapJson(peerMapJsonStr string, selfAddr string) {

}

func TestPeerListRebalance() {
	peers := NewPeerList(5, 4)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected := NewPeerList(5, 4)
	expected.Add("1111", 1)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	expected.Add("-1-1", -1)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = NewPeerList(5, 2)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected = NewPeerList(5, 2)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = NewPeerList(5, 4)
	peers.Add("1111", 1)
	peers.Add("7777", 7)
	peers.Add("9999", 9)
	peers.Add("11111111", 11)
	peers.Add("2020", 20)
	peers.Rebalance()
	expected = NewPeerList(5, 4)
	expected.Add("1111", 1)
	expected.Add("7777", 7)
	expected.Add("9999", 9)
	expected.Add("2020", 20)
	fmt.Println(reflect.DeepEqual(peers, expected))
}
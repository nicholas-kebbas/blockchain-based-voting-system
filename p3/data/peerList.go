package data

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

type PeerList struct {
	selfId int32
	peerMap map[string]int32
	maxLength int32
	mux sync.Mutex
}

func NewPeerList(id int32, maxLength int32) PeerList {}

func(peers *PeerList) Add(addr string, id int32) {}

func(peers *PeerList) Delete(addr string) {}

func(peers *PeerList) Rebalance() {}

func(peers *PeerList) Show() string {}

func(peers *PeerList) Register(id int32) {
	peers.selfId = id
	fmt.Printf("SelfId=%v\n", id)
}

func(peers *PeerList) Copy() map[string]int32 {}

func(peers *PeerList) GetSelfId() int32 {
	return peers.selfId
}

func(peers *PeerList) PeerMapToJson() (string, error) {}

func(peers *PeerList) InjectPeerMapJson(peerMapJsonStr string, selfAddr string) {}

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
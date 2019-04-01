package data

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"container/ring"
	"sync"
)

type PeerList struct {
	selfId int32
	peerMap map[string]int32 /* Maps IP Address to ID */
	maxLength int32
	mux sync.Mutex
}

type PeerMap struct {
	Addr string `json:"addr"`
	Id int32 `json:"id"`
}

type JsonPeerList struct {
	JsonRep string
}

/* Pair will hold peermap key as value, and value as key because that makes more sense */
type Pair struct {
	Key int32
	Value string
}

type PairList []Pair

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	return p[i].Key < p[j].Key
}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PairList) GetValue(i int) string {
	return p[i].Value
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
	peers.mux.Lock()
	peers.peerMap[addr] = id
	peers.mux.Unlock()
}

func(peers *PeerList) Delete(addr string) {
	peers.mux.Lock()
	delete(peers.peerMap, addr)
	peers.mux.Unlock()
}
/* Sort all peers' Id, insert SelfId, consider the list as a cycle, and choose 16 nodes at each side of SelfId.
For example, if SelfId is 10, PeerList is [7, 8, 9, 15, 16], then the closest 4 nodes are [8, 9, 15, 16].
 */
func(peers *PeerList) Rebalance() {
	peers.mux.Lock()
	peers.peerMap["myAddr"] = peers.selfId
	newPeerMap := make (map[string]int32)
	pairList := make(PairList, len(peers.peerMap))
	correctKeys := []string{}
	//newPairList := make (PairList, 2 * peers.maxLength)
	i := 0
	/* Add to list of pairs to sort, and delete also */
	for k, v := range peers.peerMap {
		pairList[i] = Pair{v, k}
		i++
	}
	halfList := int(peers.maxLength/2)
	sort.Sort(pairList)
	r := ring.New(len(pairList))
	n := r.Len()

	for i := 0; i < n; i++ {
		r.Value = pairList[i].Value
		r = r.Next()
	}

	/* Loop through ring of the keys now that they're sorted. Values are the string keys */
	for i := 0; i < n; i++ {
		if r.Value == "myAddr" {
			/* Found the selfId, so go back -halfList */
			r = r.Move(-halfList)
			for z := 0; z < halfList; z++ {
				/* Then add what we find to the PeerList */
				correctKeys = append(correctKeys, r.Value.(string))
				r = r.Next()
			}
			/* Now get the right half */
			r = r.Next()
			for z := 0; z < halfList; z++ {
				/* Then add what we find to the PeerList */
				correctKeys = append(correctKeys, r.Value.(string))
				r = r.Next()
			}

		}
		r = r.Next()
	}
	/* Add correct keys to new peer map */
	for i:= 0; i < len(correctKeys); i++ {
		// found := false
		if val, ok := peers.peerMap[correctKeys[i]]; ok {
			/* Make sure we don't add self */
			if val != peers.selfId {
				newPeerMap[correctKeys[i]] = peers.peerMap[correctKeys[i]]
			}
		}
	}
	/* Add new peermap to peerlist */
	peers.peerMap = newPeerMap
	peers.mux.Unlock()
}

/* Putting a lock here creates deadlock so don't do it */
func(peers *PeerList) Show() string {
	show, _ := peers.PeerMapToJson()
	fmt.Println(show)
	return show
}

func(peers *PeerList) Register(id int32) {
	peers.selfId = id
	fmt.Printf("SelfId=%v\n", id)
}

/* Return copy of peer list presumably */
func(peers *PeerList) Copy() map[string]int32 {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	newMap := make(map[string]int32)
	for k,v := range peers.peerMap {
		newMap[k] = v
	}
	return newMap
}

func(peers *PeerList) GetSelfId() int32 {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	return peers.selfId

}

func (peers *PeerList) GetPeerMap() map[string]int32 {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	return peers.peerMap
}

/* TODO: Fix error checking */
func(peers *PeerList) PeerMapToJson() (string, error) {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	jsonPeerMap, err := json.Marshal(peers.peerMap)
	return string(jsonPeerMap), err
}

/* Todo: Take peerMap as json String insert each entry into own peer list except for selfAddr */
func(peers *PeerList) InjectPeerMapJson(peerMapJsonStr string, selfAddr string) {
	fmt.Println("Peer Map JSON String")
	fmt.Println(peerMapJsonStr)
	var jsonMap map[string]int32
	byteRep := []byte(peerMapJsonStr)
	err := json.Unmarshal(byteRep, &jsonMap)
	if err != nil {
		fmt.Println("Error")
	}
	for k, v := range jsonMap {
		fmt.Println("Self Address")
		fmt.Println(selfAddr)
		if k != selfAddr {
			peers.Add(k, v)
		}
	}
	fmt.Println(jsonMap)

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
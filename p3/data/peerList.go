package data

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	// "strings"
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

type Pair struct {
	Key string
	Value int32
}

type PairList []Pair

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	return p[i].Value < p[j].Value
}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PairList) GetValue(i int) int32 {
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
	fmt.Println("In Rebalance")
	peers.mux.Lock()
	defer peers.mux.Unlock()
	pairList := make(PairList, len(peers.peerMap))
	newPairList := make (PairList, len(peers.peerMap))
	i := 0
	for k, v := range peers.peerMap {
		pairList[i] = Pair{k, v}
		i++
	}
	sort.Sort(pairList)
	fmt.Println(pairList)
	/* Array is now sorted. Now grab 16 to the left, and 16 to the right, if array is < 32 */
	if peers.maxLength < peers.maxLength {
		/* get to the point of selfID. Keep count to see how far in we go. */
		counter := 0
		halfList := int(peers.maxLength/2)
		for i := range pairList {
			counter ++
			/* Once we get here, count 16 back and forward */
			if  pairList.GetValue(i) == peers.selfId {
				if counter >= halfList {
					/* grab everything to the right since counter > 16 and there should be enough */
					for z := counter; z < counter+halfList; z++ {
						newPairList = append(newPairList, pairList[z])
					}
					/* grab all available to the left of the counter */
					for z := counter; z > 0; z-- {
						newPairList = append(newPairList, pairList[z])
					}
					remaining := halfList - counter
					/* Grab the remaining from the end of the array */
					for z := len(pairList); z > len(pairList) - remaining; z-- {
						newPairList = append(newPairList, pairList[z])
					}
				} else {
					for z := counter; z < len(pairList); z++ {
						newPairList = append(newPairList, pairList[z])
					}
					/* grab 16 to the left of the counter */
					for z := counter; z > counter - halfList; z-- {
						newPairList  = append(newPairList, pairList[z])
					}
					remaining := halfList - counter
					/* Grab the remaining from the beginning of the array */
					for z := 0; z > remaining; z++ {
						newPairList = append(newPairList, pairList[z])
					}
				}
				break
			}
		}
		/* Now put the new array back into the map */
	}
	fmt.Println("End of Rebalance")
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
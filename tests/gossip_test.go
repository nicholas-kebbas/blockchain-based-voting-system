package tests

import (
	"fmt"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p3/data"
	"reflect"
	"testing"
)

func TestRebalance(t *testing.T) {
	peerList := data.NewPeerList(1, 32)
	peerList.Add("address2", 8)
	peerList.Add("address0", 6)
	peerList.Add("address1", 7)
	peerList.Add("address3", 12)
	peerList.Add("address4", 15)
	peerList.Rebalance()
}

func TestRebalance2(t *testing.T) {
	peers := data.NewPeerList(5, 4)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected := data.NewPeerList(5, 4)
	expected.Add("1111", 1)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	expected.Add("-1-1", -1)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = data.NewPeerList(5, 2)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected = data.NewPeerList(5, 2)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = data.NewPeerList(5, 4)
	peers.Add("1111", 1)
	peers.Add("7777", 7)
	peers.Add("9999", 9)
	peers.Add("11111111", 11)
	peers.Add("2020", 20)
	peers.Rebalance()
	expected = data.NewPeerList(5, 4)
	expected.Add("1111", 1)
	expected.Add("7777", 7)
	expected.Add("9999", 9)
	expected.Add("2020", 20)
	fmt.Println(reflect.DeepEqual(peers, expected))
}
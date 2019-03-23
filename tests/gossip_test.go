package tests

import (
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p3/data"
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
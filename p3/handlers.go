package p3

import (
	//"../p2"
	"./data"
	"fmt"

	//"github.com/gorilla/mux"
	//"io"
	//"io/ioutil"
	//"math/rand"
	"net/http"
	//"os"
	//"strconv"
	//"strings"
	//"time"
)

var TA_SERVER = "http://localhost:6688"
var REGISTER_SERVER = TA_SERVER + "/peer"
var BC_DOWNLOAD_SERVER = TA_SERVER + "/upload"
var SELF_ADDR = "http://localhost:6686"

/* This is the canonical blockchain */
var SBC data.SyncBlockChain
var Peers data.PeerList
var ifStarted bool

func init() {
	// This function will be executed before everything else.
	// Do some initialization here.
}

// Register ID, download BlockChain, start HeartBeat
func Start(w http.ResponseWriter, r *http.Request) {}

// Display peerList and sbc
func Show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n%s", Peers.Show(), SBC.Show())
}

// Register to TA's server, get an ID
func Register() {}

// Download blockchain from TA server
func Download() {}

// Upload blockchain to whoever called this method, return jsonStr
func Upload(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := SBC.BlockChainToJson()
	if err != nil {
		// data.PrintError(err, "Upload")
	}
	fmt.Fprint(w, blockChainJson)
}

// Upload a block to whoever called this method, return jsonStr
func UploadBlock(w http.ResponseWriter, r *http.Request) {
	block := SBC.
}

/*  Received a heartbeat and follow these steps:
 Add the sender’s IP address, along with sender’s PeerList into its own PeerList
 At this time, the number of peers stored in PeerList might exceed 32 and it is ok
 If the HeartBeatData contains a new block, the node will first check if the previous block exists
 If the previous block doesn't exist, the node will ask every peer at "/block/{height}/{hash}" to download that block
 After making sure previous block exists, insert the block from HeartBeatData to the current BlockChain
*/

func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {

}

// Ask another server to return a block of certain height and hash
func AskForBlock(height int32, hash string) {}

func ForwardHeartBeat(heartBeatData data.HeartBeatData) {}

func StartHeartBeat() {}
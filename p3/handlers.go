package p3

import (
	//"../p2"
	"./data"
	"fmt"
	"io/ioutil"

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
/* Need to set up a tunnel to connect to this */
var REGISTER_SERVER = TA_SERVER + "/peer"
// First node will have the canonical block chain on it. It will first create the blockchain, and
// the other nodes will listen for it.
var FIRST_NODE = "http://localhost:6689"
var BC_DOWNLOAD_SERVER = FIRST_NODE + "/upload"
var SELF_ADDR = "http://localhost:6670"

/* This is the canonical blockchain */
var SBC data.SyncBlockChain
var Peers data.PeerList
var ifStarted bool

/*  // This function will be executed before everything else.
	// So our node should launch, then immediately grab the blockchain from BC_DOWNLOAD_SERVER
	// Start() */
func init() {

}

// Register ID, download BlockChain, start HeartBeat
func Start(w http.ResponseWriter, r *http.Request) {
	// blockChainJson, err := SBC.BlockChainToJson()
	//if err != nil {
	//	/* Report the error */
	//	// data.PrintError(err, "Upload")
	//}
	/* Register ID */
	Register()
	req, _ := http.NewRequest("GET", BC_DOWNLOAD_SERVER, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

// Display peerList and sbc
func Show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n%s", Peers.Show(), SBC.Show())
}

// Register to TA's server, get an ID
func Register() {
	req, _ := http.NewRequest("GET", REGISTER_SERVER, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

// Download blockchain from TA server
/* So launch node on SELF_ADDR and then download from the TA_SERVER
 SELF_ADDR needs to make a request to FIRST_NODE
 */
func Download() {
	/* Make http request and download from BC_DOWNLOAD_SERVER */
	// BC_DOWNLOAD_SERVER
}

// Upload blockchain to whoever called this method, return jsonStr
func Upload(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := SBC.BlockChainToJson()
	fmt.Println(r.Header)
	fmt.Println(r.Body)
	fmt.Println(r.Method)
	if err != nil {
		// data.PrintError(err, "Upload")
	}
	w.Write([]byte(blockChainJson))
	fmt.Fprint(w, blockChainJson)
}

// Upload a block to whoever called this method, return jsonStr
func UploadBlock(w http.ResponseWriter, r *http.Request) {
	// block := SBC.
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
func AskForBlock(height int32, hash string) {
	SBC.GetBlock(height, hash)
}

func ForwardHeartBeat(heartBeatData data.HeartBeatData) {

}

func StartHeartBeat() {

}
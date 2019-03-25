package p3

import (
	//"../p2"
	"./data"
	"fmt"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p1"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	//"github.com/gorilla/mux"
	//"io"
	//"math/rand"
	"net/http"
	//"os"
	//"time"
)

var TA_SERVER = "http://localhost:6688"
/* Need to set up a tunnel to connect to this */
var REGISTER_SERVER = TA_SERVER + "/peer"
// First node will have the canonical block chain on it. It will first create the blockchain, and
// the other nodes will listen for it.
var FIRST_NODE = "http://localhost:6670"
var BC_DOWNLOAD_SERVER = FIRST_NODE + "/upload"
var SELF_ADDR = "http://localhost:6671"

/* This is the canonical blockchain */
var SBC data.SyncBlockChain
var Peers data.PeerList
var ifStarted bool
var ID int32

/*  // This function will be executed before everything else.
	// So our node should launch, then immediately grab the blockchain from BC_DOWNLOAD_SERVER
	// Start()
*/
func init() {
	Register()
	go StartHeartBeat()
}

// Register ID, download BlockChain, start HeartBeat
func Start(w http.ResponseWriter, r *http.Request) {
	// blockChainJson, err := SBC.BlockChainToJson()
	//if err != nil {
	//	/* Report the error */
	//	// data.PrintError(err, "Upload")
	//}

	req, _ := http.NewRequest("GET", BC_DOWNLOAD_SERVER, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("body", string(body))
}

// Display peerList and sbc
func Show (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n%s", Peers.Show(), SBC.Show())
}

// Register to TA's server, get an ID
func Register() {
	req, _ := http.NewRequest("GET", REGISTER_SERVER, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	i, err := strconv.ParseInt(string(body), 10, 32)
	if err != nil {
		panic(err)
	}
	ID = int32(i)
}

/* Create the initial canonical Blockchain and add write it to the server */
func Create (w http.ResponseWriter, r *http.Request) {
	newBlockChain := data.NewBlockChain()
	mpt := p1.MerklePatriciaTrie{}
	mpt.Initial()
	mpt.Insert("Initial", "Value")
	newBlockChain.GenBlock(mpt)
	fmt.Println(newBlockChain)
	/* Set Global variable SBC to be this new blockchain */
	SBC = newBlockChain
	blockChainJson, _ := SBC.BlockChainToJson()
	fmt.Println(blockChainJson)
	/* Write this to the server */
	w.Write([]byte(blockChainJson))
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
	if err != nil {
		// data.PrintError(err, "Upload")
	}
	fmt.Fprint(w, blockChainJson)
}

//  If you have the block, return the JSON string of the specific block; if you don't have the block,
//  return HTTP 204: StatusNoContent; if there's an error, return HTTP 500: InternalServerError.
func UploadBlock(w http.ResponseWriter, r *http.Request) {
	url := r.RequestURI
	fmt.Println(url)
	splitURL := strings.Split(url, "/")
	blockHeightString := splitURL[2]
	blockHash := splitURL[3]
	i, err := strconv.ParseInt(blockHeightString, 10, 32)
	if err != nil {
		w.WriteHeader(500)
		panic(err)
	}
	blockHeight := int32(i)

	block, found := SBC.GetBlock(blockHeight, blockHash)
	/* Found it so write the JSONblock to output */
	if found == true {
		w.WriteHeader(200)
		fmt.Fprint(w, block.EncodeToJSON())
	} else {
		w.WriteHeader(204)
	}

}

/*  Received a heartbeat and follow these steps:
 Add the sender’s IP address, along with sender’s PeerList into its own PeerList
 At this time, the number of peers stored in PeerList might exceed 32 and it is ok
 If the HeartBeatData contains a new block, the node will first check if the previous block exists
 If the previous block doesn't exist, the node will ask every peer at "/block/{height}/{hash}" to download that block
 After making sure previous block exists, insert the block from HeartBeatData to the current BlockChain
*/

func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {
	jsonPeerList, _ := Peers.PeerMapToJson()
	data.PrepareHeartBeatData(&SBC, ID, jsonPeerList, SELF_ADDR)
	// Peers.InjectPeerMapJson()
}

// Ask another server to return a block of certain height and hash
func AskForBlock(height int32, hash string) {
	SBC.GetBlock(height, hash)
}

func ForwardHeartBeat(heartBeatData data.HeartBeatData) {

}

func StartHeartBeat() {
	for range time.Tick(time.Second *5){
		fmt.Println("Foo")
	}
}
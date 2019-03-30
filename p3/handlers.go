package p3

import (
	//"../p2"
	"./data"
	"encoding/json"
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
	//"time"
)

/* DON'T NEED TO USE THIS ANYMORE */
var TA_SERVER = "http://localhost:9000"
var REGISTER_SERVER = TA_SERVER + "/peer"
// First node will have the canonical block chain on it. It will first create the blockchain, and
// the other nodes will listen for it.
var FIRST_NODE = "http://localhost:6688"
var BC_DOWNLOAD_SERVER = FIRST_NODE + "/upload"
var SELF_ADDR string
var SELF_ADDR_FULL = "http://" + SELF_ADDR

/* This is the canonical blockchain */
var SBC data.SyncBlockChain
var Peers data.PeerList
var ifStarted bool
var ID int32

type JsonString struct {
	JsonBody string
}


/*  // This function will be executed before everything else.
	// So our node should launch, then immediately grab the blockchain from BC_DOWNLOAD_SERVER
	// Start()
*/
func init() {

}

// Register ID, download BlockChain, start HeartBeat
func Start(w http.ResponseWriter, r *http.Request) {
	// blockChainJson, err := SBC.BlockChainToJson()
	//if err != nil {
	//	/* Report the error */
	//	// data.PrintError(err, "Upload")
	//}
	/* Get address and ID */
	url := r.RequestURI
	fmt.Println(url)
	splitURL := strings.Split(url, "/")
	fmt.Println(splitURL)
	/* Get port number and set that to ID */
	fmt.Println(r.Host)
	/* Save localhost as Addr */
	fmt.Println(r.URL.Path)
	splitHostPort := strings.Split(r.Host, ":")
	i, err := strconv.ParseInt(splitHostPort[1], 10, 32)
	if err != nil {
		w.WriteHeader(500)
		panic(err)
	}
	/* ID is now port number. Address is now correct Address */
	ID = int32(i)
	SELF_ADDR = r.Host
	/* Need to instantiate the peer list */
	Peers = data.NewPeerList(ID, 32)
	Download()
	go StartHeartBeat()
}

// Display peerList and sbc
func Show (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n%s", Peers.Show(), SBC.Show())
}

// Register to TA's server, get an ID. Not needed anymore.
func Register() {
	//req, _ := http.NewRequest("GET", REGISTER_SERVER, nil)
	//res, _ := http.DefaultClient.Do(req)
	//defer res.Body.Close()
	//body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	//i, err := strconv.ParseInt(string(body), 10, 32)
	//if err != nil {
	//	panic(err)
	//}
	//ID = int32(i)
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


	/* Do most of start. Just don't download because that would be downloading from self */
	/* Get address and ID */
	url := r.RequestURI
	fmt.Println(url)
	splitURL := strings.Split(url, "/")
	fmt.Println(splitURL)
	/* Get port number and set that to ID */
	fmt.Println(r.Host)
	/* Save localhost as Addr */
	fmt.Println(r.URL.Path)
	splitHostPort := strings.Split(r.Host, ":")
	i, err := strconv.ParseInt(splitHostPort[1], 10, 32)
	if err != nil {
		w.WriteHeader(500)
		panic(err)
	}
	/* ID is now port number. Address is now correct Address */
	ID = int32(i)
	SELF_ADDR = r.Host
	/* Need to instantiate the peer list */
	Peers = data.NewPeerList(ID, 32)
}

// Download blockchain from First Node
/* So launch node on SELF_ADDR and then download from the TA_SERVER
 SELF_ADDR needs to make a request to FIRST_NODE
 */
func Download() {
	/* Make http request and download from BC_DOWNLOAD_SERVER */
	// So Download makes a POST request to First Node Server/upload.
	// Then the response will be HeartBeatdata (configure that in Upload())
	//
	// BC_DOWNLOAD_SERVER
	jsonPeerList, _ := Peers.PeerMapToJson()
	newHeartBeatData := data.PrepareHeartBeatData(&SBC, ID, jsonPeerList, SELF_ADDR)
	fmt.Println("Peer Map JSON")
	fmt.Println(newHeartBeatData.HeartBeatToJson())
	/* Need to figure out what to send here in the request */
	res, _ := http.Post(BC_DOWNLOAD_SERVER, "application/json; charset=UTF-8", strings.NewReader(newHeartBeatData.HeartBeatToJson()))
	//res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

// Called By Download (POST Request from local node's Download Function)
func Upload(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := SBC.BlockChainToJson()
	if err != nil {
		// data.PrintError(err, "Upload")
	}
	/* Also store the ID and Address of the incoming request */
	splitHostPort := strings.Split(r.Host, ":")
	i, err := strconv.ParseInt(splitHostPort[1], 10, 32)
	if err != nil {
		w.WriteHeader(500)
		panic(err)
	}
	/* ID is now port number. Address is now correct Address */
	body, err := ioutil.ReadAll(r.Body)
	remoteId := int32(i)
	Peers.Add(r.Host, remoteId)
	fmt.Println("remoteId")
	fmt.Println(remoteId)
	/* Send POST request to /upload with Address and ID data. Then populate the peer list */
	fmt.Fprint(w, blockChainJson)
	s := string(body)
	fmt.Println(s)
	var t data.HeartBeatData
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	fmt.Println(t.Addr)
	fmt.Println(t.Id)
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
 Add the remote address, and the PeerMapJSON into local PeerMap.
 At this time, the number of peers stored in PeerList might exceed 32 and it is ok
 Then check if the HeartBeatData contains a new block.
 If so, do these: (1) check if the parent block exists.
 If not, call AskForBlock() to download the parent block.
(2) insert the new block from HeartBeatData.
(3) HeartBeatData.hops minus one, and if it's still bigger than 0,
call ForwardHeartBeat() to forward this heartBeat to all peers.
 If the HeartBeatData contains a new block, the node will first check if the previous block exists
 If the previous block doesn't exist, the node will ask every peer at "/block/{height}/{hash}" to download that block
 After making sure previous block exists, insert the block from HeartBeatData to the current BlockChain
*/

func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {
	/* Parse the request */
	fmt.Println("Request Host")
	fmt.Println(r.Host)
	// Peers.Add(r.Host, r.)
	jsonPeerList, _ := Peers.PeerMapToJson()
	data.PrepareHeartBeatData(&SBC, ID, jsonPeerList, SELF_ADDR)

	// Peers.InjectPeerMapJson()
}

// Ask another server to return a block of certain height and hash
func AskForBlock(height int32, hash string) {
	SBC.GetBlock(height, hash)
}

/* Send to all the peers. Will probably want to send a post request to their ReceiveHeartBeat */
func ForwardHeartBeat(heartBeatData data.HeartBeatData) {

	/* Get the peerMap */
	localPeerMap := Peers.GetPeerMap()
	for k,v := range localPeerMap {
		remoteAddress := k
		fmt.Println(remoteAddress)
		fmt.Println(v)
		// resp, err := http.Post(REGISTER_SERVER, "application/json; charset=UTF-8", strings.NewReader(string(heartBeatData)))
	}
}

func StartHeartBeat() {
	for range time.Tick(time.Second *5){
		/* PrepareHeartBeatData() to create a HeartBeatData,
		and send it to all peers in the local PeerMap */
		stringJson, _ := Peers.PeerMapToJson()
		newHeartBeatData := data.PrepareHeartBeatData(&SBC, ID, stringJson, SELF_ADDR)
		ForwardHeartBeat(newHeartBeatData)
		fmt.Println(ID)
		fmt.Println(SELF_ADDR)
	}
}
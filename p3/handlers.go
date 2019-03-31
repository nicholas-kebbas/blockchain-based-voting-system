package p3

import (
	"./data"
	"encoding/json"
	"fmt"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p1"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p2"
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
var FIRST_NODE_ADDRESS = "localhost:6688"
var FIRST_NODE_PORT int32 = 6688
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
	if len(Peers.GetPeerMap()) == 0 {
		Peers = data.NewPeerList(ID, 32)
	}

	/* Need to also add the first node that we're connecting to, as long as this isn't Node 1 */
	if SELF_ADDR != FIRST_NODE_ADDRESS {
		Peers.Add(FIRST_NODE_ADDRESS, FIRST_NODE_PORT)
		Download()
	}
	fmt.Println("Starting Heartbeat in", SELF_ADDR)
	go StartHeartBeat()
}

// Display peerList and sbc
func Show (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the Peer List: ")
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
	/* Instantiate and grab the blockchain */
	fmt.Println("Body in Download")
	SBC = data.NewBlockChain()
	SBC.UpdateEntireBlockChain(string(body))
	fmt.Println("String body")
	fmt.Println(string(body))
	fmt.Println("Final Updated Blockchain")
	fmt.Println(SBC)
	fmt.Println("As JSON")
	fmt.Println(SBC.BlockChainToJson())
}

// Called By Download (POST Request from local node's Download Function)
func Upload(w http.ResponseWriter, r *http.Request) {

	/* Also store the ID and Address of the incoming request */
	splitHostPort := strings.Split(r.Host, ":")
	i, err := strconv.ParseInt(splitHostPort[1], 10, 32)
	if err != nil {
		fmt.Println(i)
		w.WriteHeader(500)
		panic(err)
	}

	/* ID is now port number. Address is now correct Address */
	body, err := ioutil.ReadAll(r.Body)
	/* Send POST request to /upload with Address and ID data. Then populate the peer list */
	s := string(body)
	fmt.Println(s)
	var t data.HeartBeatData
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	Peers.Add(t.Addr, t.Id)
	fmt.Println(t.Addr)
	fmt.Println(t.Id)
	/* Response should be the block chain and the peer list */
	blockChainJson, err := SBC.BlockChainToJson()
	if err != nil {
		// data.PrintError(err, "Upload")
	}
	fmt.Println("Blockchain in Upload Function")
	fmt.Println(blockChainJson)
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
 Add the remote address, and the PeerMapJSON into local PeerMap.
 At this time, the number of peers stored in PeerList might exceed 32 and it is ok
 0. Then check if the HeartBeatData contains a new block.
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
	fmt.Println("Receiving Heartbeat")
	fmt.Println(r.Host)
	// Peers.Add(r.Host, r.)
	jsonPeerList, _ := Peers.PeerMapToJson()
	data.PrepareHeartBeatData(&SBC, ID, jsonPeerList, SELF_ADDR)

	/* Parse the request and add peers to this node's peer map */
	body, err := ioutil.ReadAll(r.Body)
	/* Send POST request to /upload with Address and ID data. Then populate the peer list */
	s := string(body)
	fmt.Println(s)
	var t data.HeartBeatData
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	/* Adding Sender to peer list as well as Peers */
	Peers.Add(t.Addr, t.Id)
	Peers.InjectPeerMapJson(t.PeerMapJson ,SELF_ADDR)
	/* Check if heartbeat contains new block */
	if t.IfNewBlock == true {
		/* If it contains new block, check if our blockchain has parent */
		newBlock := p2.Block{}
		newBlock.DecodeFromJson(t.BlockJson)
		/* Might need to lock here */
		if SBC.CheckParentHash(newBlock) == false {
			/* We do -1 because we need the parents height */
			AskForBlock(newBlock.Header.Height - 1, newBlock.Header.ParentHash)
			fmt.Println("Adding parent hash, next line should always be true")
		}
		if SBC.CheckParentHash(newBlock) == true {
			fmt.Println("Inserting new Block")
			SBC.Insert(newBlock)
			t.Hops = t.Hops - 1
			/* If still greater than 1, forward on */
			if t.Hops > 0 {
				ForwardHeartBeat(t)
			}
		}
	}
	fmt.Println("End of receive heartbeat")
	// Peers.InjectPeerMapJson()
}

// Ask another server to return a block of certain height and hash
func AskForBlock(height int32, hash string) {
	fmt.Println("Asking for block")
	localPeerMap := Peers.GetPeerMap()
	for k, _ := range localPeerMap {
		remoteAddress := "http://" + k + "/block/" + string(height) + "/" + hash
		/* Make a GET request to peer to see if the block is there */
		req, _ := http.NewRequest("GET", remoteAddress, nil)
		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()
		if res.StatusCode == 204 {
			body, _ := ioutil.ReadAll(res.Body)
			fmt.Println(string(body))
			/* May not need to do any of this unmarshaling. Leaving here for now */
			newJsonBlock := p2.JsonBlock{}
			err := json.Unmarshal(body, newJsonBlock)
			if err != nil {
				fmt.Println("Could not unmarshal block")
			}
			newBlock := p2.Block{}
			newBlock.DecodeFromJson(string(body))
			SBC.Insert(newBlock)
			break
		}
		res.Body.Close()

	}
	SBC.GetBlock(height, hash)
}

/* Send to all the peers. Will probably want to send a post request to their ReceiveHeartBeat */
func ForwardHeartBeat(heartBeatData data.HeartBeatData) {
	fmt.Println("Forward Heart Beat. Breaking on rebalance")
	/* Need to rebalance before send. Makes the most sense to do it here */
	Peers.Rebalance()
	fmt.Println("Rebalanced")
	/* Get the peerMap */
	localPeerMap := Peers.GetPeerMap()
	for k,v := range localPeerMap {
		remoteAddress := "http://" + k + "/heartbeat/receive"
		fmt.Println(remoteAddress)
		fmt.Println(v)
		fmt.Println("Data sending")
		fmt.Println(heartBeatData.HeartBeatToJson())
		resp, _ := http.Post(remoteAddress, "application/json; charset=UTF-8", strings.NewReader(heartBeatData.HeartBeatToJson()))
		fmt.Println(resp)
	}
}

func StartHeartBeat() {
	for range time.Tick(time.Second *5) {
		fmt.Println("Heartbeat")
		/* PrepareHeartBeatData() to create a HeartBeatData,
		and send it to all peers in the local PeerMap */
		stringJson, _ := Peers.PeerMapToJson()
		newHeartBeatData := data.PrepareHeartBeatData(&SBC, ID, stringJson, SELF_ADDR)
		ForwardHeartBeat(newHeartBeatData)
	}
}
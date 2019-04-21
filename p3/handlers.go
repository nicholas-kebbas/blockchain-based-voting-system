package p3

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p1"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p2"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p3/data"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/voting"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"math/rand"
	"net/http"
)

// First node will have the canonical block chain on it. It will first create the blockchain, and
// the other nodes will listen for it.
var FIRST_NODE_ADDRESS = "localhost:6688"
var FIRST_NODE_PORT int32 = 6688
var FIRST_NODE = "http://localhost:6688"
var FIRST_NODE_SERVER = FIRST_NODE + "/upload"
var FIRST_NODE_BALLOT = FIRST_NODE + "/ballot"
var SELF_ADDR string
var FOUNDREMOTE = false
var CREATED = false
/* Adding permissioning to blockchain */

/* For simplicity's sake, this can just be the Port Numbers since we're collecting that info already.

In production, we can actually keep a seperate list of predetermined allowed IDs
 */

 /* Only these IDs are allowed to write */
var ALLOWED_IDS = map[int32]bool{
	6688:true,
	6669:true,
	6670:true,
}

/* Need to check this ID whenever write attempts are made

So upon create
 */

/* This is the canonical blockchain */
var SBC data.SyncBlockChain
var Peers data.PeerList
var ifStarted bool
var ID int32
var BALLOT voting.Ballot

type JsonString struct {
	JsonBody string
}

type CanonicalChain struct {
	blocks []CanonicalChainBlock
}

type CanonicalChainBlock struct {
	height int32
	timestamp int64
	hash string
	parentHash string
	size int32
}


/*  // This function will be executed before everything else.
	// So our node should launch, then immediately grab the blockchain from BC_DOWNLOAD_SERVER
	// Start()
*/
func init() {
	ifStarted = false
}

// Register ID, download BlockChain, start HeartBeat
func Start(w http.ResponseWriter, r *http.Request) {
	/* Check if node is in set of Allowed IDs */
	if _, ok := ALLOWED_IDS[ID]; ok {
		/* Get address and ID */
		if ifStarted != true {
			splitHostPort := strings.Split(r.Host, ":")
			i, err := strconv.ParseInt(splitHostPort[1], 10, 32)
			if err != nil {
				w.WriteHeader(500)
				panic(err)
			}
			/* ID is now port number. Address is now correct Address */
			ID = int32(i)
			SELF_ADDR = r.Host
			RandSeed()
			/* Need to instantiate the peer list */
			if len(Peers.Copy()) == 0 {
				Peers = data.NewPeerList(ID, 32)
			}

			/* Need to also add the first node that we're connecting to, as long as this isn't Node 1 */
			if SELF_ADDR != FIRST_NODE_ADDRESS {
				Peers.Add(FIRST_NODE_ADDRESS, FIRST_NODE_PORT)
				Download()
			}
			/* Call down the voting ballot */
			go StartHeartBeat()
			go StartTryingNonces()
			ifStarted = true
		}
	}
}

// Display peerList and sbc
func Show (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the Peer List: ")
	fmt.Fprintf(w, "%s\n%s", Peers.Show(), SBC.Show())
}

/* Create the initial canonical Blockchain and add write it to the server */
func Create (w http.ResponseWriter, r *http.Request) {
	/* This is an SBC */
	if CREATED == false {

		/* Move the checking of ID up first to confirm this is allowed */
		/* Do most of start. Just don't download because that would be downloading from self */
		/* Get address and ID */
		/* Get port number and set that to ID */
		/* Save localhost as Addr */
		splitHostPort := strings.Split(r.Host, ":")
		i, err := strconv.ParseInt(splitHostPort[1], 10, 32)
		if err != nil {
			w.WriteHeader(500)
			panic(err)
		}
		/* ID is now port number. Address is now correct Address */
		ID = int32(i)
		SELF_ADDR = r.Host
		/* Check if ID is allowed in ALLOWED_IDs */
		if _, ok := ALLOWED_IDS[ID]; ok {
			newBlockChain := data.NewBlockChain()
			mpt := p1.MerklePatriciaTrie{}
			mpt.Initial()
			mpt.Insert("Initial", "Value")
			newBlockChain.GenBlock(mpt)
			/* Set Global variable SBC to be this new blockchain */
			SBC = newBlockChain
			blockChainJson, _ := SBC.BlockChainToJson()
			/* Write this to the server */
			w.Write([]byte(blockChainJson))

			/* Need to instantiate the peer list */
			Peers = data.NewPeerList(ID, 32)
			CREATED = true
		}
	}
}

// Download blockchain from First Node
/* So launch node on SELF_ADDR and then download from the TA_SERVER
 SELF_ADDR needs to make a request to FIRST_NODE
 */
func Download() {
	/* Make http request and download from BC_DOWNLOAD_SERVER */
	// So Download makes a POST request to First Node Server/upload.
	// Then the response will be HeartBeatdata (configure that in Upload())
	jsonPeerList, _ := Peers.PeerMapToJson()
	/* Just creating trie here so we can use the prepareHeartBeatData Function */
	trie := p1.MerklePatriciaTrie{}
	newHeartBeatData := data.PrepareHeartBeatData(&SBC, ID, jsonPeerList, SELF_ADDR, false, "", trie)
	/* Need to figure out what to send here in the request */
	res, _ := http.Post(FIRST_NODE_SERVER, "application/json; charset=UTF-8", strings.NewReader(newHeartBeatData.HeartBeatToJson()))
	//res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	/* Instantiate and grab the blockchain */
	SBC = data.NewBlockChain()
	SBC.UpdateEntireBlockChain(string(body))
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
	/* Response should be the block chain and the peer list */
	blockChainJson, err := SBC.BlockChainToJson()
	if err != nil {
		// data.PrintError(err, "Upload")
	}
	fmt.Fprint(w, blockChainJson)
}

/* Download the uploaded ballot from other Nodes */
func DownloadBallot() {
	/* Just creating trie here so we can use the prepareHeartBeatData Function */
	/* Need to figure out what to send here in the request */
	res, _ := http.Post(FIRST_NODE_BALLOT, "application/json; charset=UTF-8", strings.NewReader(newHeartBeatData.HeartBeatToJson()))
	//res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	/* Instantiate and grab the blockchain */
	SBC.UpdateEntireBlockChain(string(body))
}

/* POST the contents of the ballot so other nodes can download. Read from ballot.json */
func UploadBallot(w http.ResponseWriter, r *http.Request) {
	/* Read the JSON File */
	plan, _ := ioutil.ReadFile("ballot.json")
	var data interface{}
	err := json.Unmarshal(plan, &data)
	if err != nil {
		fmt.Println("Error")
	}
	fmt.Fprint(w, plan)
	fmt.Println(plan)
}

//  If you have the block, return the JSON string of the specific block; if you don't have the block,
//  return HTTP 204: StatusNoContent; if there's an error, return HTTP 500: InternalServerError.
func UploadBlock(w http.ResponseWriter, r *http.Request) {
	url := r.RequestURI
	splitURL := strings.Split(url, "/")
	blockHeightString := splitURL[2]
	cleanHeight := strings.Replace(blockHeightString, "%0", "", 1)
	cleanHeight = strings.Replace(blockHeightString, "%", "", 1)
	blockHash := splitURL[3]
	i, err := strconv.ParseInt(cleanHeight, 10, 32)
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

/* Pretty Print the canonical chain */

/*  If there's only one block at height 100, the chain of that block and all its predecessors
(parent, parent of the parent, etc) is the canonical chain. */


/* If there are multiple blocks at max height n, they are considered as forks, and each fork can form a chain.
The canonical chain would be decided once one of the forks grows and that chain becomes the longest chain.  */
func Canonical(w http.ResponseWriter, r *http.Request) {
	/* First get latest blocks. If there are more than 1 here, that will mean there are
	two equally long forks.
	 */
	var i = 0
	latestBlocks := SBC.GetLatestBlocks()
	canonicalChains := []CanonicalChain{}
	for i = 0; i < len(latestBlocks); i++ {
		canonicalChain := CanonicalChain{}
		latestBlock := latestBlocks[i]
		canonicalChain.writeInfoToCanonicalChain(latestBlock)
		canonicalChains = append(canonicalChains, canonicalChain)
	}

	/* Then loop through, print their information or record to a data structure,
	and do the same for parent.
	*/
	output := ""
	var z = 0
	for i = 0; i < len(canonicalChains); i++ {
		fmt.Fprint(w, "Chain: ")
		fmt.Fprintln(w, strconv.Itoa(i + 1))
		for z = 0; z < len(canonicalChains[i].blocks); z++ {
			output = fmt.Sprintf("%+v", canonicalChains[i].blocks[z])
			fmt.Fprintln(w, output)
		}
	}

	/* Need this at end or else golang/browser strips out line breaks */
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

}

func (canonicalChain *CanonicalChain) writeInfoToCanonicalChain(block p2.Block) {
	if block.Header.Height == 1 {
		chainBlock := CanonicalChainBlock{}
		chainBlock.height = block.Header.Height
		chainBlock.timestamp = block.Header.TimeStamp
		chainBlock.hash = block.Header.Hash
		chainBlock.parentHash = block.Header.ParentHash
		chainBlock.size = block.Header.Size
		canonicalChain.blocks = append(canonicalChain.blocks, chainBlock)
		return
	}
	chainBlock := CanonicalChainBlock{}
	chainBlock.height = block.Header.Height
	chainBlock.timestamp = block.Header.TimeStamp
	chainBlock.hash = block.Header.Hash
	chainBlock.parentHash = block.Header.ParentHash
	chainBlock.size = block.Header.Size
	canonicalChain.blocks = append(canonicalChain.blocks, chainBlock)
	parentBlock,_ := SBC.GetParentBlock(block)
	canonicalChain.writeInfoToCanonicalChain(parentBlock)
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

	/* Parse the request and add peers to this node's peer map */
	body, err := ioutil.ReadAll(r.Body)
	/* Send POST request to /upload with Address and ID data. Then populate the peer list */
	var t data.HeartBeatData
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	/* Adding Sender to peer list as well as Sender's Peerlist */
	Peers.Add(t.Addr, t.Id)
	Peers.InjectPeerMapJson(t.PeerMapJson, SELF_ADDR)
	/* So now, this is if we received a block and need to verify the nonce */
	if t.IfNewBlock == true {
		/* If it contains new block, check if our blockchain has parent */
		newBlock := p2.Block{}
		newBlock = newBlock.DecodeFromJson(t.BlockJson)
		/* First verify the nonce, then do the other steps */
		if VerifyNonceFromBlock(newBlock) {
			FOUNDREMOTE = true
			/* This keeps returning false when it shouldn't. We have the parent in most cases, so this func is wrong */
			if SBC.CheckParentHash(newBlock) == false {
				/* So we're looking for the parent here and adding it if we don't have it.*/
				AskForBlock(newBlock.Header.Height - 1, newBlock.Header.ParentHash)
			}
			if SBC.CheckParentHash(newBlock) == true {
				SBC.Insert(newBlock)
				t.Hops = t.Hops - 1
				/* If still greater than 1, forward on */
				if t.Hops > 0 {
					ForwardHeartBeat(t)
				}
			}
			/* So now we may or may not have found it */
			FOUNDREMOTE = false
			/* Stop our own search for nonce and start it again with new parent */

		}
	}
}

/* Function recursively asks for all the missing predecessor blocks instead of only the parent block.  */
func AskForBlock(height int32, hash string) {
	localPeerMap := Peers.Copy()
	if height == 0 {
		fmt.Print("Parent Block not in chain")
		return
	}
	for k, _ := range localPeerMap {
		s := strconv.FormatInt(int64(height), 10)
		/* Calls Upload Block */
		remoteAddress := "http://" + k + "/block/" + s + "/" + hash
		/* Make a GET request to peer to see if the block is there */
		req, _ := http.NewRequest("GET", remoteAddress, nil)
		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()
		/* This means no block */
		if res.StatusCode == 204 {
			/* Recurse and Look for Parent Block if we can't find parent */
			parentHeight := height - 1
			AskForBlock(parentHeight, hash)
		} else if res.StatusCode == 200 {
			body, _ := ioutil.ReadAll(res.Body)
			newBlock := p2.Block{}
			/* Decode from Json must not be working correctly */
			newBlock.DecodeFromJson(string(body))
			SBC.Insert(newBlock)
			break
		} else {
			fmt.Println("500 error")
		}
		res.Body.Close()
	}
}

/* Send to all the peers. Will probably want to send a post request to their ReceiveHeartBeat */
func ForwardHeartBeat(heartBeatData data.HeartBeatData) {
	/* Need to rebalance before send. Makes the most sense to do it here */
	Peers.Rebalance()
	/* Get the peerMap */
	localPeerMap := Peers.Copy()
	for k, _ := range localPeerMap {
		remoteAddress := "http://" + k + "/heartbeat/receive"
		http.Post(remoteAddress, "application/json; charset=UTF-8", strings.NewReader(heartBeatData.HeartBeatToJson()))
	}
}

func StartHeartBeat() {
	for range time.Tick(time.Second *7) {
		/* PrepareHeartBeatData() to create a HeartBeatData,
		and send it to all peers in the local PeerMap */
		stringJson, _ := Peers.PeerMapToJson()
		trie := p1.MerklePatriciaTrie{}
		newHeartBeatData := data.PrepareHeartBeatData(&SBC, ID, stringJson, SELF_ADDR, false, "", trie)
		ForwardHeartBeat(newHeartBeatData)
	}
}

func SendBlock(nonce string, trie p1.MerklePatriciaTrie) {
	fmt.Println("New Block Found! Sending")
	/* PrepareHeartBeatData() to create a HeartBeatData,
	and send it to all peers in the local PeerMap */
	stringJson, _ := Peers.PeerMapToJson()
	newHeartBeatData := data.PrepareHeartBeatData(&SBC, ID, stringJson, SELF_ADDR, true, nonce, trie)
	ForwardHeartBeat(newHeartBeatData)
}

/*
This function starts a new thread that tries different nonces to generate new blocks.
Nonce is a string of 16 hexes such as "1f7b169c846f218a".
Initialize the rand when you start a new node with something unique about each node, such as the current time or the port number.

 */
func StartTryingNonces() {
	/* Start an outer while loop. This will need to run the whole time the node is active. */
	for ok := true; ok; ok = true {
		/* Get the latest block or one of the latest blocks to use as a parent block. */
		parentHash := ""
		if SBC.GetLength() == 0 {
			parentHash = "Genesis"
		} else {
			parentBlock := SBC.GetLatestBlocks()[0]
			parentHash = parentBlock.Header.Hash
		}
		/* Create an MPT. */
		mpt := GenerateRandomMpt()

		/*  Randomly generate the first nonce, verify it with simple PoW algorithm to see if
			SHA3(parentHash + nonce + mptRootHash) starts with 10 0's (or the number you modified into).
			Since we use one laptop to try different nonces, six to seven 0's could be enough.
			If the nonce failed the verification, increment it by 1 and try the next nonce.
			 */
		var i= 1
		for i < 1000 {
			if FOUNDREMOTE == false {
				nonce := p2.CalculateNonce()
				str := parentHash + nonce + mpt.GetRoot()
				hash := sha3.Sum256([]byte(str))
				/* Get first N digits of sha */
				firstN := string(hash[:])
				if strings.HasPrefix(firstN, "000") {
					/* Create block, etc. */
					i = 1000
					/* Send the block to peers */
					SendBlock(nonce, mpt)
					/* Break out of loop and start working on next nonce */
					break
				}
			} else {
				/* Was found remotely so try next nonce now */
				break
			}
		}

		/* If a nonce is found and the next block is generated, forward that block to all peers with a HeartBeatData;
		  In case we want to do that outside of loop */

		/*  If someone else found a nonce first, and you received the new block through your function ReceiveHeartBeat(),
		  stop trying nonce on the current block, continue to the while loop by jumping to the step(2)
		   */
	} /* End outer while. This will never end unless node is stopped. */
}


func VerifyNonceFromBlock(block p2.Block) bool {
	verified := false
	str := block.Header.ParentHash + block.Header.Nonce + block.GetMptRoot()
	hash := sha3.Sum256([]byte(str))
	/* Get first N digits of sha */
	firstN := string(hash[:])
	if strings.HasPrefix(firstN, "000") {
		verified = true
	}
	return verified
}

func GenerateRandomMpt() p1.MerklePatriciaTrie {
	mpt := p1.MerklePatriciaTrie{}
	mpt.Initial()
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	for text != "q" {  // break the loop if text == "q"
		fmt.Print("Enter your Vote: ")
		scanner.Scan()
		text = scanner.Text()
		if text != "q" {
			fmt.Println("You voted for ", text)
		}
	}
	fmt.Println("Thanks for voting!")
	/* Record the Vote in the MPT. Value can be vote value, but need to decide what the key is */
	mpt.Insert(text, text)

	if scanner.Err() != nil {
		// handle error.
	}
	return mpt
}

func RandSeed() {
	rand.Seed(int64(ID))
}